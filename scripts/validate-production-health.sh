#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"
TIMEOUT_SECONDS="${FTN_HEALTH_TIMEOUT_SECONDS:-120}"
INTERVAL_SECONDS="${FTN_HEALTH_INTERVAL_SECONDS:-2}"

log(){ printf '[FTN][HEALTH] %s\n' "$*"; }
fail(){ printf '[FTN][HEALTH][ERROR] %s\n' "$*" >&2; exit 1; }
trap 'printf "[FTN][HEALTH][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

[ "$(id -u)" -eq 0 ] || fail 'Run as root'
command -v docker >/dev/null 2>&1 || fail 'Docker is required'
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'
docker compose version >/dev/null 2>&1 || fail 'Docker Compose v2 is required'
[ -f "$ENV_FILE" ] || fail ".env is missing: $ENV_FILE"
[[ "$TIMEOUT_SECONDS" =~ ^[0-9]+$ ]] && (( TIMEOUT_SECONDS > 0 )) || fail 'FTN_HEALTH_TIMEOUT_SECONDS must be a positive integer'
[[ "$INTERVAL_SECONDS" =~ ^[0-9]+$ ]] && (( INTERVAL_SECONDS > 0 )) || fail 'FTN_HEALTH_INTERVAL_SECONDS must be a positive integer'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest is marked for FTN live deployment'

check_stack(){
  local compose="$1"
  local elapsed=0
  local expected actual service state health exit_code restart
  local config_json
  config_json="$(docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" config --format json)"
  mapfile -t expected < <(printf '%s' "$config_json" | python3 -c 'import json,sys; d=json.load(sys.stdin); print(*d.get("services",{}).keys(),sep="\n")')
  ((${#expected[@]})) || fail "No services resolved from $compose"

  while (( elapsed <= TIMEOUT_SECONDS )); do
    local ps
    ps="$(docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" ps -a --format '{{.Service}}|{{.State}}|{{.Health}}|{{.ExitCode}}' 2>/dev/null || true)"
    local all_ok=1
    local missing=''
    for service in "${expected[@]}"; do
      actual="$(printf '%s\n' "$ps" | awk -F'|' -v s="$service" '$1==s {print; exit}')"
      if [ -z "$actual" ]; then
        all_ok=0
        missing+=" $service"
        continue
      fi
      IFS='|' read -r _service state health exit_code <<< "$actual"
      restart="$(printf '%s' "$config_json" | python3 -c 'import json,sys; d=json.load(sys.stdin); s=sys.argv[1]; v=d.get("services",{}).get(s,{}).get("restart"); print("" if v is None else v)' "$service")"
      if [ "$restart" = "no" ]; then
        if [ "$state" != "exited" ] || [ "$exit_code" != "0" ]; then
          all_ok=0
          log "${compose#$ROOT_DIR/}: $service one-shot state=$state exit_code=$exit_code"
        fi
        continue
      fi
      if [ "$state" != "running" ]; then
        all_ok=0
        log "${compose#$ROOT_DIR/}: $service state=$state"
        continue
      fi
      if [ -n "$health" ] && [ "$health" != "healthy" ]; then
        all_ok=0
        log "${compose#$ROOT_DIR/}: $service health=$health"
      fi
    done
    if (( all_ok )); then return 0; fi
    sleep "$INTERVAL_SECONDS"
    elapsed=$((elapsed + INTERVAL_SECONDS))
  done

  log "FAILED: ${compose#$ROOT_DIR/}"
  docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" ps >&2 || true
  docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" logs --tail=100 >&2 || true
  [ -z "$missing" ] || log "Missing service(s):$missing"
  return 1
}

for compose in "${manifests[@]}"; do
  log "Checking: ${compose#$ROOT_DIR/}"
  check_stack "$compose" || fail "Production health gate failed: ${compose#$ROOT_DIR/}"
done

log 'Production service health: PASS'
