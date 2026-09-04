#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROFILE="$ROOT_DIR/configs/system/99-ftn-routing.conf"
SYSCTL_DIR="/etc/sysctl.d"
TARGET="$SYSCTL_DIR/99-ftn-routing.conf"

log(){ printf '[FTN][ROUTING-KERNEL] %s\n' "$*"; }
fail(){ printf '[FTN][ROUTING-KERNEL][ERROR] %s\n' "$*" >&2; exit 1; }
trap 'printf "[FTN][ROUTING-KERNEL][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

[ "$(id -u)" -eq 0 ] || fail 'Run as root'
command -v sysctl >/dev/null 2>&1 || fail 'sysctl is required'
[ -f "$PROFILE" ] || fail "routing kernel profile missing: $PROFILE"

install -d -m 0755 "$SYSCTL_DIR"
install -m 0644 "$PROFILE" "$TARGET"
sysctl --system >/dev/null

require_value(){
  local key="$1" expected="$2" actual
  actual="$(sysctl -n "$key")"
  [ "$actual" = "$expected" ] || fail "$key=$actual; expected $expected"
}

require_value net.ipv4.conf.all.rp_filter 2
require_value net.ipv4.conf.default.rp_filter 2
require_value net.ipv4.ip_forward 1
require_value net.ipv6.conf.all.forwarding 1

log 'FTN routing kernel profile applied and validated'
