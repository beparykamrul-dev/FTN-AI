#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/frontend/control-center"

echo "[1/4] Installing frontend dependencies..."
npm install --quiet

echo "[2/4] Building TypeScript..."
npm run build --if-present

echo "[3/4] Starting FTN Control Center development server..."
echo "✓ Frontend available at http://localhost:3000"
npm run dev

