#!/bin/sh
set -eu
BASE="${FTN_BASE_URL:-http://127.0.0.1:8080}"
for path in /healthz /readyz /api/v1/services; do
  code="$(curl -sS -o /dev/null -w '%{http_code}' "$BASE$path")"
  [ "$code" = "200" ] || { echo "FAIL $path ($code)"; exit 1; }
done
echo "FTN control-plane production checks: PASS"
