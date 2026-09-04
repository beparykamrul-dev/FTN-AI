#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"
source "$ROOT_DIR/scripts/production-compose-env.sh"

fail(){ printf '[FTN][STORAGE][ERROR] %s\n' "$*" >&2; exit 1; }
log(){ printf '[FTN][STORAGE] %s\n' "$*"; }
trap 'printf "[FTN][STORAGE][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

command -v docker >/dev/null 2>&1 || fail 'Docker is required'
command -v python3 >/dev/null 2>&1 || fail 'python3 is required'
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'
[ -f "$ENV_FILE" ] || fail '.env is missing'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest found'

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

for compose in "${manifests[@]}"; do
  docker compose --env-file "$ENV_FILE" -f "$compose" config --format json \
    | python3 -c 'import json,sys
x=json.load(sys.stdin)
for service,spec in x.get("services",{}).items():
    for mount in spec.get("volumes",[]) or []:
        if not isinstance(mount,dict): continue
        typ=mount.get("type"); source=mount.get("source"); target=mount.get("target")
        if typ in ("bind","volume") and source and target: print(f"{typ}\t{source}\t{target}\t{service}")' >> "$tmp"
done

while IFS=$'\t' read -r typ source target service; do
  [ -n "$source" ] || continue
  if [ "$typ" = bind ]; then
    case "$source" in /*) ;; *) fail "Unresolved relative bind source: $source ($service -> $target)" ;; esac
    [ -e "$source" ] || log "Bind path will be created by Compose: $source"
  fi
done < "$tmp"

awk -F '\t' '$1=="bind" {key=$2; target=$3; if (seen[key] && seen_target[key] != target) {print "bind source mounted at incompatible targets: " key " => " seen_target[key] " and " target; bad=1} seen[key]=1; seen_target[key]=target} END{exit bad?1:0}' "$tmp" || fail 'Conflicting bind-mount ownership detected'
awk -F '\t' '$1=="volume" {key=$2; target=$3; if (seen[key] && seen_target[key] != target) {print "named volume mounted at incompatible targets: " key " => " seen_target[key] " and " target; bad=1} seen[key]=1; seen_target[key]=target} END{exit bad?1:0}' "$tmp" || fail 'Conflicting named-volume ownership detected'

log 'Production storage ownership: PASS'
