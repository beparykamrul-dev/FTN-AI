#!/usr/bin/env bash
# Canonical production Compose profile selection.
# Unprofiled services are always live; optional production services are enabled
# explicitly. The wireless/Kismet profile remains opt-in because it requires a
# host wireless interface and elevated capabilities.
FTN_PRODUCTION_PROFILES="${FTN_PRODUCTION_PROFILES:-observability,backup,enterprise,acme}"
export COMPOSE_PROFILES="$FTN_PRODUCTION_PROFILES"
