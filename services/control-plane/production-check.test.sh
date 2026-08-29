#!/bin/sh
set -eu
command -v sh >/dev/null
command -v curl >/dev/null
test -x production-check.sh || chmod +x production-check.sh
grep -q '/healthz' production-check.sh
grep -q '/readyz' production-check.sh
grep -q '/api/v1/services' production-check.sh
echo 'production-check script: PASS'
