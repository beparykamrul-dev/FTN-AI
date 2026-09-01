#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "========================================="
echo "FTN-AI Development Environment"
echo "========================================="
echo ""

# Start backend in background
echo "Starting Backend Diagnostics Service..."
bash "$ROOT/scripts/dev/start-backend.sh" &
BACKEND_PID=$!

sleep 2

# Start frontend
echo ""
echo "Starting Frontend Control Center..."
bash "$ROOT/scripts/dev/start-frontend.sh" &
FRONTEND_PID=$!

echo ""
echo "========================================="
echo "✓ FTN-AI Development Environment Ready"
echo "========================================="
echo ""
echo "Frontend:  http://localhost:3000"
echo "Backend:   http://localhost:8080"
echo "Dashboard: http://localhost:3000/overview"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

wait

