#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

command -v go >/dev/null && gofmt -l services/control-plane/*.go >/tmp/ftn-gofmt.txt || true
if [[ -s /tmp/ftn-gofmt.txt ]]; then
  cat /tmp/ftn-gofmt.txt
  exit 1
fi

GOPROXY=direct go test ./...
if [[ -f services/control-plane/go.mod ]]; then
  (cd services/control-plane && GOPROXY=direct go test ./...)
fi

echo "FTN final integrity check: PASS"
