#!/usr/bin/env bash
set -Eeuo pipefail

fail(){ printf '[FTN][ROUTING-KERNEL][ERROR] %s\n' "$*" >&2; exit 1; }

command -v sysctl >/dev/null 2>&1 || fail 'sysctl is required'

check(){
  local key="$1" expected="$2" actual
  actual="$(sysctl -n "$key" 2>/dev/null || true)"
  [ "$actual" = "$expected" ] || fail "$key=$actual; expected $expected"
}

check net.ipv4.conf.all.rp_filter 2
check net.ipv4.conf.default.rp_filter 2
check net.ipv4.ip_forward 1
check net.ipv6.conf.all.forwarding 1

# Linux applies the effective rp_filter as the maximum of all and the
# interface-specific value. Validate every currently-present IPv4 interface so
# provider-facing links cannot silently run with strict/disabled RPF settings.
for path in /proc/sys/net/ipv4/conf/*/rp_filter; do
  [ -f "$path" ] || continue
  iface="${path%/rp_filter}"
  iface="${iface##*/}"
  actual="$(cat "$path")"
  [ "$actual" = 2 ] || fail "interface $iface rp_filter=$actual; expected 2"
done

printf '[FTN][ROUTING-KERNEL] PASS: rp_filter=2 and forwarding are enforced\n'
