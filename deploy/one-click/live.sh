#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

log(){ printf '[FTN][LIVE] %s\n' "$*"; }
fail(){ printf '[FTN][LIVE][ERROR] %s\n' "$*" >&2; exit 1; }
trap 'printf "[FTN][LIVE][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

[ "$(id -u)" -eq 0 ] || fail 'Run as root: sudo bash deploy/one-click/live.sh'
command -v docker >/dev/null 2>&1 || fail 'Docker is required; run bootstrap.sh first'
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

ENV_FILE="$ROOT_DIR/.env"
[ -f "$ENV_FILE" ] || fail '.env is missing; run bootstrap.sh first'
chmod 600 "$ENV_FILE"

# Only manifests explicitly marked production are started. This prevents test/dev
# compose files from accidentally entering a live deployment.
mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | xargs -0 -r grep -Il 'FTN_PRODUCTION_STACK=true' | sort)

((${#manifests[@]})) || fail 'No production Compose manifest is marked with FTN_PRODUCTION_STACK=true'

log "Found ${#manifests[@]} production stack(s)"
for compose in "${manifests[@]}"; do
  log "Validate: ${compose#$ROOT_DIR/}"
  docker compose --env-file "$ENV_FILE" -f "$compose" config --quiet
 done

# Run repository-level preflight before mutating live services.
if [ -x "$ROOT_DIR/scripts/preflight.sh" ]; then
  log 'Running repository preflight'
  "$ROOT_DIR/scripts/preflight.sh"
fi

# Apply database migrations through the service's declared migration runner.
if [ -f "$ROOT_DIR/services/control-plane/docker-compose.yml" ]; then
  log 'Starting control-plane database/runtime foundation'
  docker compose --env-file "$ENV_FILE" -f "$ROOT_DIR/services/control-plane/docker-compose.yml" up -d --build --remove-orphans
fi

# Start every other explicitly production-marked stack.
for compose in "${manifests[@]}"; do
  [ "$compose" = "$ROOT_DIR/services/control-plane/docker-compose.yml" ] && continue
  log "Start: ${compose#$ROOT_DIR/}"
  docker compose --env-file "$ENV_FILE" -f "$compose" up -d --build --remove-orphans
 done

# Wait for the control-plane readiness contract when it exists.
if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1 || \
   curl -fsS --max-time 3 http://127.0.0.1:8080/healthz >/dev/null 2>&1; then
  log 'Control-plane readiness: OK'
elif [ -f "$ROOT_DIR/services/control-plane/docker-compose.yml" ]; then
  for _ in $(seq 1 60); do
    if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1 || \
       curl -fsS --max-time 3 http://127.0.0.1:8080/healthz >/dev/null 2>&1; then
      break
    fi
    sleep 2
  done
  curl -fsS --max-time 5 http://127.0.0.1:8080/readyz >/dev/null 2>&1 || \
    curl -fsS --max-time 5 http://127.0.0.1:8080/healthz >/dev/null 2>&1 || \
    fail 'Control-plane readiness check failed'
  log 'Control-plane readiness: OK'
fi

log 'Production stack status'
for compose in "${manifests[@]}"; do
  docker compose --env-file "$ENV_FILE" -f "$compose" ps
 done

log 'ONE-CLICK LIVE DEPLOYMENT COMPLETE'
