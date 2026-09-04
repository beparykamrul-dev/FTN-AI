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

set -a
# shellcheck disable=SC1091
. "$RELEASE_ROOT/.env"
set +a
log 'Checking database schema compatibility before rollback'
bash "$previous/deploy/one-click/release-compatibility.sh" "$previous"
log 'Database downgrade: NEVER (schema remains append-only)'

mapfile -t manifests < <(find "$previous" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production manifests in previous release'

# Validate everything before changing the release pointer or restarting services.
for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$compose" config --quiet
  log "Validated: ${compose#$previous/}"
done

# Reject host-port collisions in the rollback target before touching the live stack.
tmp_ports="$(mktemp)"
trap 'rm -f "$tmp_ports"' EXIT
for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$compose" config --format json \
    | python3 -c 'import json,sys
x=json.load(sys.stdin)
for s,v in x.get("services",{}).items():
  for p in v.get("ports",[]) or []:
    if isinstance(p,dict) and p.get("published") is not None:
      print(f"{p.get(\"published\")}/{p.get(\"protocol\",\"tcp\")} {s}")' >> "$tmp_ports"
 done
awk '{ key=$1; if (seen[key]++) { print "duplicate host port: " $0; bad=1 } } END { exit bad ? 1 : 0 }' "$tmp_ports" \
  || fail 'Rollback target has a host port collision'
rm -f "$tmp_ports"
trap - ERR

ln -sfn "$previous" "$CURRENT_LINK"
cd "$previous"

CONTROL_COMPOSE="$previous/services/control-plane/docker-compose.yml"
log 'Starting previous PostgreSQL foundation'
docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$CONTROL_COMPOSE" up -d postgres

pg_ok=0
for _ in $(seq 1 60); do
  health="$(docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$CONTROL_COMPOSE" ps -a --format '{{.Service}} {{.Health}}' 2>/dev/null || true)"
  if printf '%s\n' "$health" | awk '$1=="postgres" && $2=="healthy" {ok=1} END{exit ok?0:1}'; then pg_ok=1; break; fi
  sleep 2
done
[ "$pg_ok" -eq 1 ] || fail 'PostgreSQL did not become healthy during rollback'

log 'Starting previous control-plane without database downgrade'
docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" \
  -f "$CONTROL_COMPOSE" up -d --build --remove-orphans control-plane

ready=0
for _ in $(seq 1 60); do
  if curl -fsS --max-time 3 http://127.0.0.1:8080/readyz >/dev/null 2>&1; then ready=1; break; fi
  sleep 2
done
[ "$ready" -eq 1 ] || {
  docker compose --profile '*' --env-file "$RELEASE_ROOT/.env" -f "$CONTROL_COMPOSE" logs --tail=200 control-plane postgres migration-runner >&2 || true
  fail 'Previous release failed control-plane readiness after rollback'
}

# Restore every production-marked stack after the canonical control-plane is ready.
# The shared ftn-control-plane network is therefore available before Control Center.
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

log 'Application rollback complete; all production stacks restored; database schema was not downgraded'
