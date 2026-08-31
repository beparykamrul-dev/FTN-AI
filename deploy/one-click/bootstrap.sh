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

# Runtime credentials are generated locally and never committed.
# Repair an incomplete/malformed .env automatically so a previous bad value
# cannot block the one-click deployment.
touch .env
chmod 600 .env

DB_PASSWORD=''
if grep -q '^FTN_DB_PASSWORD=' .env; then
  DB_PASSWORD="$(sed -n 's/^FTN_DB_PASSWORD=//p' .env | tail -n 1)"
fi

if [ -z "$DB_PASSWORD" ]; then
  DB_PASSWORD="$(openssl rand -hex 32)"
  tmp_env="$(mktemp)"
  awk -v p="$DB_PASSWORD" '
    BEGIN{done=0}
    /^FTN_DB_PASSWORD=/{if(!done){print "FTN_DB_PASSWORD=" p; done=1}; next}
    {print}
    END{if(!done) print "FTN_DB_PASSWORD=" p}
  ' .env > "$tmp_env"
  chmod 600 "$tmp_env"
  mv "$tmp_env" .env
  log 'Generated/ repaired FTN_DB_PASSWORD in local runtime .env'
fi

# Never print the credential. Verify only that a non-empty value exists.
grep -q '^FTN_DB_PASSWORD=.' .env || fail 'Unable to create FTN_DB_PASSWORD'

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
  docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" config >/dev/null
  log 'Starting FTN stack'
  docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" up -d --build --remove-orphans

  log 'Waiting for control-plane health'
  healthy=0
  for _ in $(seq 1 30); do
    if curl -fsS http://127.0.0.1:8080/healthz >/dev/null 2>&1; then
      healthy=1
      break
    fi
    sleep 2
  done

  if [ "$healthy" -eq 1 ]; then
    log 'Control-plane health: OK'
  else
    log 'Control-plane health check did not become ready; showing service state'
    docker compose --env-file "$ROOT_DIR/.env" -f "$COMPOSE_FILE" ps
    exit 1
  fi
else
  fail 'No Compose manifest found; cannot start the FTN stack'
fi

log 'FTN one-click bootstrap completed successfully'
