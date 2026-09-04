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
[ -L "$CURRENT_LINK" ] || fail "$CURRENT_LINK is not a release symlink"

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
[ -f "$previous/deploy/one-click/live.sh" ] || fail 'Previous release lacks canonical live runner'

# Never attempt an automatic database downgrade. Rollback only changes the application
# release pointer, then starts the previous release against the existing schema.
log "Current release: $current"
log "Previous release: $previous"
log 'Database downgrade: NEVER (schema remains append-only)'

ln -sfn "$previous" "$CURRENT_LINK"
cd "$previous"

[ -f "$RELEASE_ROOT/.env" ] || fail 'Production .env is missing'
chmod 600 "$RELEASE_ROOT/.env"

log 'Validating previous release Compose configuration'
mapfile -t manifests < <(find "$previous" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production manifests in previous release'
for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$compose" config --quiet
done

log 'Starting previous application release without database rollback'
docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" \
  -f "$previous/services/control-plane/docker-compose.yml" up -d --build --remove-orphans control-plane

ready=0
for _ in $(seq 1 60); do
  if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done
[ "$ready" -eq 1 ] || fail 'Previous release failed readiness after rollback'

log 'Application rollback complete; database schema was not downgraded'
