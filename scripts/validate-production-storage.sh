#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"

fail(){ printf '[FTN][STORAGE][ERROR] %s\n' "$*" >&2; exit 1; }
log(){ printf '[FTN][STORAGE] %s\n' "$*"; }
trap 'printf "[FTN][STORAGE][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

command -v docker >/dev/null 2>&1 || fail 'Docker is required'
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'
[ -f "$ENV_FILE" ] || fail '.env is missing'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest found'

TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" config --format json \
    | python3 -c 'import json,sys
x=json.load(sys.stdin)
for service, spec in x.get("services",{}).items():
    for mount in spec.get("volumes",[]) or []:
        if not isinstance(mount,dict):
            continue
        typ=mount.get("type")
        source=mount.get("source")
        target=mount.get("target")
        if typ in ("bind","volume") and source and target:
            print(f"{typ}\t{source}\t{target}\t{service}")' >> "$TMP"
done

# A bind source is a host path; reject relative paths that Compose has not
# resolved and reject the same host source being mounted at incompatible
# targets across independent production stacks.
while IFS=$'\t' read -r typ source target service; do
  [ -n "$source" ] || continue
  if [ "$typ" = bind ]; then
    case "$source" in
      /*) ;;
      *) fail "Unresolved relative bind source: $source ($service -> $target)" ;;
    esac
    [ -e "$source" ] || log "Bind path will be created by Compose: $source"
  fi
done < "$TMP"

awk -F '\t' '
$1=="bind" { key=$2; target=$3; owner=$4; if (seen[key] && (seen_target[key] != target)) { print "bind source mounted at incompatible targets: " key " => " seen_target[key] " and " target; bad=1 } seen[key]=owner; seen_target[key]=target }
END { exit bad ? 1 : 0 }
' "$TMP" || fail 'Conflicting bind-mount ownership detected'

# Named volume names are global to Docker, even when Compose projects differ.
# A volume reused by multiple stacks is only safe when it serves the same
# container target; otherwise one stack can mutate another stack's data.
awk -F '\t' '$1=="volume" { key=$2; target=$3; if (seen[key] && seen_target[key] != target) { print "named volume mounted at incompatible targets: " key " => " seen_target[key] " and " target; bad=1 } seen[key]=1; seen_target[key]=target } END { exit bad ? 1 : 0 }' "$TMP" \
  || fail 'Conflicting named-volume ownership detected'

log 'Production storage ownership: PASS'
