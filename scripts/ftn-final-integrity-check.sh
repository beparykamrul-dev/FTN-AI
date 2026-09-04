#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

mapfile -t tracked_go < <(git ls-files -- '*.go')
: > /tmp/ftn-gofmt.txt
if command -v go >/dev/null 2>&1 && ((${#tracked_go[@]})); then
  gofmt -l "${tracked_go[@]}" > /tmp/ftn-gofmt.txt
fi
if [[ -s /tmp/ftn-gofmt.txt ]]; then
  cat /tmp/ftn-gofmt.txt
  exit 1
fi

GOPROXY=direct go test ./...
if [[ -f services/control-plane/go.mod ]]; then
  (cd services/control-plane && GOPROXY=direct go test ./...)
fi

echo "FTN final integrity check: PASS"
