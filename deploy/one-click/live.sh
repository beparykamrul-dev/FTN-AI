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
command -v curl >/dev/null 2>&1 || fail 'curl is required'
command -v python3 >/dev/null 2>&1 || fail 'python3 is required'
ENV_FILE="$ROOT_DIR/.env"
[ -f "$ENV_FILE" ] || fail '.env is missing; run bootstrap.sh first'
chmod 600 "$ENV_FILE"
mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest is marked for FTN live deployment'
[ -f "$ROOT_DIR/scripts/preflight.sh" ] || fail 'Production preflight script is missing'
log 'Running repository preflight'
bash "$ROOT_DIR/scripts/preflight.sh"
log "Found ${#manifests[@]} production stack(s)"
for compose in "${manifests[@]}"; do
  log "Validate: ${compose#$ROOT_DIR/}"
  docker compose --profile "*" --env-file "$ENV_FILE" -f "$compose" config --quiet
done
[ -f "$ROOT_DIR/scripts/validate-production-storage.sh" ] || fail 'Production storage validator is missing'
[ -f "$ROOT_DIR/scripts/validate-production-ports.sh" ] || fail 'Production port validator is missing'
log 'Validating persistent storage ownership'
bash "$ROOT_DIR/scripts/validate-production-storage.sh"
log 'Validating production host-port ownership'
bash "$ROOT_DIR/scripts/validate-production-ports.sh"
CONTROL_COMPOSE="$ROOT_DIR/services/control-plane/docker-compose.yml"
if [ -f "$CONTROL_COMPOSE" ]; then
  mapfile -t control_services < <(docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" config --services)
  for required in postgres migration-runner control-plane; do
    printf '%s\n' "${control_services[@]}" | grep -Fxq "$required" || fail "control-plane Compose service missing: $required"
  done
  log 'Starting PostgreSQL foundation'
  docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" up -d postgres
  log 'Waiting for PostgreSQL health'
  for _ in $(seq 1 60); do
    health="$(docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" ps -a --format '{{.Service}} {{.Health}}' 2>/dev/null || true)"
    if printf '%s\n' "$health" | awk '$1=="postgres" && $2=="healthy" {ok=1} END{exit ok?0:1}'; then break; fi
    sleep 2
  done
  health="$(docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" ps -a --format '{{.Service}} {{.Health}}' 2>/dev/null || true)"
  printf '%s\n' "$health" | awk '$1=="postgres" && $2=="healthy" {ok=1} END{exit ok?0:1}' || fail 'PostgreSQL did not become healthy'
  log 'Applying database migrations (fail-closed)'
  docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" up --force-recreate --abort-on-container-exit --exit-code-from migration-runner migration-runner
  log 'Starting control-plane application'
  docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" up -d --build --remove-orphans control-plane
fi
for compose in "${manifests[@]}"; do
  [ "$compose" = "$CONTROL_COMPOSE" ] && continue
  log "Start: ${compose#$ROOT_DIR/}"
  docker compose --profile "*" --env-file "$ENV_FILE" -f "$compose" up -d --build --remove-orphans
done
if [ -f "$CONTROL_COMPOSE" ]; then
  ready=0
  for _ in $(seq 1 60); do
    if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1; then ready=1; break; fi
    sleep 2
  done
  [ "$ready" -eq 1 ] || { docker compose --profile "*" --env-file "$ENV_FILE" -f "$CONTROL_COMPOSE" logs --tail=200 control-plane postgres migration-runner >&2 || true; fail 'Control-plane readiness check failed'; }
  log 'Control-plane readiness: OK'
fi
log 'Production stack status'
for compose in "${manifests[@]}"; do docker compose --profile "*" --env-file "$ENV_FILE" -f "$compose" ps; done
log 'ONE-CLICK LIVE DEPLOYMENT COMPLETE'
