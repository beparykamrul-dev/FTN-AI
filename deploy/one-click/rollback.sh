#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RELEASE_ROOT="${FTN_RELEASE_ROOT:-/opt/ftn-ai}"
CURRENT_LINK="$RELEASE_ROOT/current"

log(){ printf '[FTN][ROLLBACK] %s\n' "$*"; }
fail(){ printf '[FTN][ROLLBACK][ERROR] %s\n' "$*" >&2; exit 1; }
trap 'printf "[FTN][ROLLBACK][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

[ "$(id -u)" -eq 0 ] || fail 'Run as root: sudo bash deploy/one-click/rollback.sh'
command -v docker >/dev/null 2>&1 || fail 'Docker is required'
command -v readlink >/dev/null 2>&1 || fail 'readlink is required'
command -v curl >/dev/null 2>&1 || fail 'curl is required'
[ -L "$CURRENT_LINK" ] || fail "$CURRENT_LINK is not a release symlink"
[ -f "$RELEASE_ROOT/.env" ] || fail 'Production .env is missing'

current="$(readlink -f "$CURRENT_LINK")"
[ -d "$current" ] || fail "Current release does not exist: $current"

previous=""
mapfile -t releases < <(find "$RELEASE_ROOT/releases" -mindepth 1 -maxdepth 1 -type d -printf '%T@ %p\n' 2>/dev/null | sort -nr | cut -d' ' -f2-)
for release in "${releases[@]}"; do
  if [ "$release" != "$current" ] && [ -f "$release/deploy/one-click/live.sh" ]; then
    previous="$release"
    break
  fi
done
[ -n "$previous" ] || fail 'No previous deployable release found'

[ -f "$previous/services/control-plane/docker-compose.yml" ] || fail 'Previous release lacks control-plane Compose manifest'
[ -f "$previous/deploy/one-click/release-compatibility.sh" ] || fail 'Previous release lacks schema compatibility gate'

chmod 600 "$RELEASE_ROOT/.env"
log "Current release: $current"
log "Previous release: $previous"
log 'Checking database schema compatibility before rollback'
set -a
# shellcheck disable=SC1091
. "$RELEASE_ROOT/.env"
set +a
bash "$previous/deploy/one-click/release-compatibility.sh" "$previous"

# Never attempt an automatic database downgrade. Rollback changes application code
# and production containers only; the schema remains append-only.
log 'Database downgrade: NEVER (schema remains append-only)'
ln -sfn "$previous" "$CURRENT_LINK"
cd "$previous"

mapfile -t manifests < <(find "$previous" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production manifests in previous release'

log 'Validating previous release Compose configuration'
for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$compose" config --quiet
done

CONTROL_COMPOSE="$previous/services/control-plane/docker-compose.yml"
log 'Starting previous control-plane release without database downgrade'
docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" \
  -f "$CONTROL_COMPOSE" up -d --build --remove-orphans control-plane

ready=0
for _ in $(seq 1 60); do
  if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done
[ "$ready" -eq 1 ] || {
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$CONTROL_COMPOSE" logs --tail=200 control-plane postgres migration-runner >&2 || true
  fail 'Previous release failed readiness after rollback'
}

# Restore every other production-marked stack as well. Starting only control-plane
# leaves DNS, observability, SFU and the Control Center on the failed release.
for compose in "${manifests[@]}"; do
  [ "$compose" = "$CONTROL_COMPOSE" ] && continue
  log "Restoring: ${compose#$previous/}"
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" \
    -f "$compose" up -d --build --remove-orphans
 done

log 'Verifying restored production stack state'
for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$compose" ps
 done

log 'Application rollback complete; database schema was not downgraded'
