#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

log(){ printf '[FTN] %s\n' "$*"; }
fail(){ printf '[FTN][ERROR] %s\n' "$*" >&2; exit 1; }

[ "$(id -u)" -eq 0 ] || fail 'Run as root: sudo bash deploy/one-click/bootstrap.sh'
export DEBIAN_FRONTEND=noninteractive

if command -v apt-get >/dev/null 2>&1; then
  log 'Installing required host dependencies'
  apt-get update
  apt-get install -y ca-certificates curl git jq openssl postgresql-client procps
fi

if ! command -v docker >/dev/null 2>&1; then
  log 'Installing Docker Engine'
  curl -fsSL https://get.docker.com | sh
fi

systemctl enable --now docker 2>/dev/null || true
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

mkdir -p data logs backups secrets
chmod 700 secrets

touch .env
chmod 600 .env

ensure_secret(){
  local key="$1" value="$2" tmp
  tmp="$(mktemp)"
  awk -v k="$key" -v v="$value" 'BEGIN{done=0} $0 ~ "^"k"=" {if(!done){print k"="v;done=1};next} {print} END{if(!done)print k"="v}' .env > "$tmp"
  chmod 600 "$tmp"
  mv "$tmp" .env
}

secret_or_generate(){
  local key="$1" value
  value="$(sed -n "s/^${key}=//p" .env | tail -n 1)"
  [ -n "$value" ] || value="$(openssl rand -hex 32)"
  ensure_secret "$key" "$value"
}

secret_or_generate FTN_DB_PASSWORD
secret_or_generate FTN_API_AUTH_TOKEN
secret_or_generate FTN_SFU_API_KEY
secret_or_generate FTN_SFU_API_SECRET
secret_or_generate FTN_TURN_PASSWORD
ensure_secret FTN_TURN_USERNAME "ftn"

log 'Applying bounded adaptive TCP performance profile'
bash "$ROOT_DIR/scripts/configure-tcp-performance.sh"

# All production Compose manifests are discovered and started by one canonical
# runner. This keeps bootstrap, upgrades and recovery on the same execution path.
chmod +x deploy/one-click/live.sh
exec bash deploy/one-click/live.sh
