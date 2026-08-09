#!/usr/bin/env bash
set -Eeuo pipefail

# FTN Server Socket Agent
# Read-only discovery/telemetry bridge for servers managed by FTN-AI.
# It exposes a local Unix socket with JSON-line requests/responses.
# Never executes arbitrary commands received over the socket.
#
# Request examples:
#   {"action":"health"}
#   {"action":"services"}
#   {"action":"info"}
#
# Start:
#   ./scripts/ftn-server-socket-agent.sh /run/ftn/server.sock

SOCKET="${1:-/run/ftn/server.sock}"
SOCKET_DIR="$(dirname "$SOCKET")"
mkdir -p "$SOCKET_DIR"
rm -f "$SOCKET"

command -v socat >/dev/null 2>&1 || {
  echo "socat is required" >&2
  exit 127
}

json_escape() {
  python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
}

respond() {
  local req="$1"
  case "$req" in
    *'"action":"health"'*)
      printf '%s\n' '{"ok":true,"agent":"ftn-server-socket-agent","mode":"read-only"}' ;;
    *'"action":"info"'*)
      printf '{"ok":true,"hostname":%s,"kernel":%s,"uptime":%s}\n' \
        "$(hostname | json_escape)" \
        "$(uname -sr | json_escape)" \
        "$(uptime -p | json_escape)" ;;
    *'"action":"services"'*)
      services="$(systemctl list-units --type=service --state=running --no-legend --no-pager 2>/dev/null | awk '{print $1}' | head -200)"
      printf '{"ok":true,"services":['
      first=1
      while IFS= read -r s; do
        [[ -z "$s" ]] && continue
        [[ $first -eq 0 ]] && printf ','
        printf '%s' "$(printf '%s' "$s" | json_escape)"
        first=0
      done <<< "$services"
      printf ']}\n' ;;
    *)
      printf '%s\n' '{"ok":false,"error":"unsupported_action"}' ;;
  esac
}

# socat creates one short-lived handler per connection.
socat "UNIX-LISTEN:$SOCKET,fork,mode=0660" SYSTEM:'read req; respond "$req"' 2>/dev/null
