#!/usr/bin/env bash
# Canonical production Compose profile selection.
# Unprofiled services are always live. The wireless/Kismet profile remains
# opt-in because it requires a host wireless interface and elevated capabilities.
SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_SOURCE="${ENV_FILE:-$SCRIPT_ROOT/.env}"
if [ -z "${FTN_PRODUCTION_PROFILES:-}" ] && [ -f "$ENV_SOURCE" ]; then
  FTN_PRODUCTION_PROFILES="$(sed -n 's/^FTN_PRODUCTION_PROFILES=//p' "$ENV_SOURCE" | tail -n 1)"
fi
FTN_PRODUCTION_PROFILES="${FTN_PRODUCTION_PROFILES:-observability,backup,enterprise,acme}"
case ",$FTN_PRODUCTION_PROFILES," in
  *,wireless,*) : ;; # explicitly permitted; host capability checks remain external
esac
export COMPOSE_PROFILES="$FTN_PRODUCTION_PROFILES"
