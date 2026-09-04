#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

log(){ printf '[FTN][PREFLIGHT] %s\n' "$*"; }
fail(){ printf '[FTN][PREFLIGHT][ERROR] %s\n' "$*" >&2; exit 1; }
trap 'printf "[FTN][PREFLIGHT][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

[ "$(id -u)" -eq 0 ] || fail 'Run as root'
command -v docker >/dev/null 2>&1 || fail 'Docker is required'
docker info >/dev/null 2>&1 || fail 'Docker daemon is unavailable'
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'

ENV_FILE="$ROOT_DIR/.env"
[ -f "$ENV_FILE" ] || fail '.env is missing; run bootstrap first'
[ "$(stat -c '%a' "$ENV_FILE" 2>/dev/null || stat -f '%Lp' "$ENV_FILE")" = '600' ] || fail '.env must have mode 0600'

[ -f "$ROOT_DIR/services/control-plane/docker-compose.yml" ] || fail 'control-plane production manifest is missing'
grep -Eq '^x-ftn-production-stack:[[:space:]]*true[[:space:]]*$' "$ROOT_DIR/services/control-plane/docker-compose.yml" || fail 'control-plane is not marked production'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest found'

for compose in "${manifests[@]}"; do
  log "Compose syntax/config: ${compose#$ROOT_DIR/}"
  docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" config --quiet
  if grep -Eiq '(password|secret|token|api[_-]?key)[[:space:]]*:[[:space:]]*(["'"']?)([^$][^[:space:]"'"']+)' "$compose"; then
    fail "possible hard-coded credential in ${compose#$ROOT_DIR/}; use environment substitution"
  fi
done

# Required migration/deployment primitives must exist before any live mutation.
for path in \
  services/control-plane/scripts/migrate.sh \
  services/control-plane/migrations/024_execution_immutability_backfill.sql \
  services/control-plane/tests/sql/024_execution_integrity_state_machine.sql \
  deploy/one-click/live.sh; do
  [ -f "$ROOT_DIR/$path" ] || fail "required deployment artifact missing: $path"
done

log "Production manifests: ${#manifests[@]}"
log 'Preflight PASS'
