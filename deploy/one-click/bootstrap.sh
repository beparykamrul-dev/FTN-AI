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
  apt-get install -y ca-certificates curl git jq openssl postgresql-client
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
  if ! grep -q "^${key}=" .env; then
    printf '%s=%s\n' "$key" "$value" >> .env
  fi
  tmp="$(mktemp)"
  awk -v k="$key" -v v="$value" 'BEGIN{done=0} $0 ~ "^"k"=" {if(!done){print k"="v;done=1};next} {print} END{if(!done)print k"="v}' .env > "$tmp"
  chmod 600 "$tmp"
  mv "$tmp" .env
}

DB_PASSWORD="$(sed -n 's/^FTN_DB_PASSWORD=//p' .env | tail -n 1)"
[ -n "$DB_PASSWORD" ] || DB_PASSWORD="$(openssl rand -hex 32)"
ensure_secret FTN_DB_PASSWORD "$DB_PASSWORD"

API_TOKEN="$(sed -n 's/^FTN_API_AUTH_TOKEN=//p' .env | tail -n 1)"
[ -n "$API_TOKEN" ] || API_TOKEN="$(openssl rand -hex 32)"
ensure_secret FTN_API_AUTH_TOKEN "$API_TOKEN"

COMPOSE_FILE=""
for candidate in \
  "services/control-plane/docker-compose.yml" \
  "docker-compose.yml" \
  "compose.yml"; do
  if [ -f "$candidate" ]; then COMPOSE_FILE="$candidate"; break; fi
done

[ -n "$COMPOSE_FILE" ] || fail 'No supported Compose manifest found'

log "Validating Compose: $COMPOSE_FILE"
docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" config --quiet

log 'Building and starting FTN control plane'
docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" up -d --build --remove-orphans

log 'Waiting for FTN control plane readiness'
ready=0
for _ in $(seq 1 60); do
  if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1 || \
     curl -fsS --max-time 3 http://127.0.0.1:8080/healthz >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done

if [ "$ready" -ne 1 ]; then
  docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" ps
  docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" logs --tail=120 control-plane postgres migration-runner || true
  fail 'FTN control plane did not become ready'
fi

log 'FTN control plane: READY'
log 'Installed at: '"$ROOT_DIR"
log 'Local endpoint: http://127.0.0.1:8080'
log 'Use the generated .env for runtime secrets; credentials are never printed.'
