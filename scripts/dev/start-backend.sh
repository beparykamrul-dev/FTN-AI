#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/backend"

echo "[1/4] Setting up Go environment..."
export GOPROXY=direct
export GOSUMDB=off

echo "[2/4] Running Go module diagnostics service..."
echo "✓ Backend diagnostics service starting on :8080"
go run ./internal/diagnostics/service.go

