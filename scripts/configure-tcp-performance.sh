#!/usr/bin/env bash
set -Eeuo pipefail

log(){ printf '[FTN][TCP] %s\n' "$*"; }
fail(){ printf '[FTN][TCP][ERROR] %s\n' "$*" >&2; exit 1; }

[ "$(id -u)" -eq 0 ] || fail 'Run as root'
command -v sysctl >/dev/null 2>&1 || fail 'sysctl is required'

# TCP_NODELAY is enabled by Go's net/http TCP listener path for accepted TCP
# connections. Do not try to emulate it with a global sysctl: Linux has no
# safe global TCP_NODELAY switch.
#
# For buffers, use Linux autotuning with a bounded ceiling. This is deliberately
# moderate: it permits high-BDP links to grow beyond the small defaults without
# allocating the ceiling to every connection.
cat >/etc/sysctl.d/99-ftn-tcp-performance.conf <<'EOF'
# FTN adaptive TCP receive/send buffering.
net.ipv4.tcp_moderate_rcvbuf = 1
net.core.rmem_max = 4194304
net.core.wmem_max = 4194304
net.ipv4.tcp_rmem = 4096 131072 4194304
net.ipv4.tcp_wmem = 4096 16384 4194304
EOF

sysctl --system >/dev/null

for key in net.ipv4.tcp_moderate_rcvbuf net.core.rmem_max net.core.wmem_max net.ipv4.tcp_rmem net.ipv4.tcp_wmem; do
  value="$(sysctl -n "$key")"
  [ -n "$value" ] || fail "Unable to read $key"
  log "$key=$value"
done

log 'TCP performance tuning: PASS (TCP_NODELAY via Go net/http; adaptive bounded buffers enabled)'
