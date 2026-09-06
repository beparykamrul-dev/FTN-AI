#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

mapfile -t tracked_go < <(git ls-files -- '*.go')
: > /tmp/ftn-gofmt.txt
if ((${#tracked_go[@]})); then
  gofmt -l "${tracked_go[@]}" > /tmp/ftn-gofmt.txt
fi
if [[ -s /tmp/ftn-gofmt.txt ]]; then
  cat /tmp/ftn-gofmt.txt
  exit 1
fi

mapfile -t modules < <(git ls-files 'go.mod' | sort -u)
if ((${#modules[@]} == 0)); then
  echo 'no Go modules found' >&2
  exit 1
fi
for mod in "${modules[@]}"; do
  dir="$(dirname "$mod")"
  echo "testing module: $dir"
  (cd "$dir" && GOPROXY=direct go test ./...)
  echo "building module: $dir"
  (cd "$dir" && GOPROXY=direct go build -trimpath ./...)
done

echo "FTN final integrity check: PASS"
