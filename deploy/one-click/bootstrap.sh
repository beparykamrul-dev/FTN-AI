#!/usr/bin/env bash
set -Eeuo pipefail

# FTN one-click bootstrap: installs only required host packages, prepares directories,
# validates configuration, and starts the compose stack when available.
# Never stores provider/payment secrets in the repository.

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
  log 'Docker not found; installing Docker Engine from the official convenience script'
  curl -fsSL https://get.docker.com | sh
fi

systemctl enable --now docker 2>/dev/null || true

docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

mkdir -p data logs backups secrets
chmod 700 secrets

[ -f .env ] || {
  if [ -f .env.example ]; then cp .env.example .env; chmod 600 .env; fi
}

if [ -f docker-compose.yml ] || [ -f compose.yml ]; then
  log 'Validating compose configuration'
  docker compose config >/dev/null
  log 'Starting FTN stack'
  docker compose up -d --remove-orphans
else
  log 'No compose manifest found; host bootstrap completed without starting services'
fi

log 'FTN bootstrap completed'
