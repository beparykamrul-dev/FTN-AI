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
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

mkdir -p data logs backups secrets
chmod 700 secrets

# Generate a local database password for first boot when the deployment has not supplied one.
# This file is runtime state and must never be committed to Git.
touch .env
chmod 600 .env
if ! grep -q '^FTN_DB_PASSWORD=' .env; then
  printf 'FTN_DB_PASSWORD=%s\n' "$(openssl rand -hex 32)" >> .env
fi

COMPOSE_FILE=''
if [ -f docker-compose.yml ]; then
  COMPOSE_FILE='docker-compose.yml'
elif [ -f compose.yml ]; then
  COMPOSE_FILE='compose.yml'
elif [ -f services/control-plane/docker-compose.yml ]; then
  COMPOSE_FILE='services/control-plane/docker-compose.yml'
fi

if [ -n "$COMPOSE_FILE" ]; then
  log "Validating Compose: $COMPOSE_FILE"
  docker compose -f "$COMPOSE_FILE" config >/dev/null
  log "Starting FTN stack"
  docker compose -f "$COMPOSE_FILE" up -d --build --remove-orphans
else
  log 'No Compose manifest found; host bootstrap completed without starting services'
fi

log 'FTN bootstrap completed'
