#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"

fail(){ printf '[FTN][PORTS][ERROR] %s\n' "$*" >&2; exit 1; }
log(){ printf '[FTN][PORTS] %s\n' "$*"; }
trap 'printf "[FTN][PORTS][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

command -v docker >/dev/null 2>&1 || fail 'Docker is required'
command -v python3 >/dev/null 2>&1 || fail 'python3 is required'
[ -f "$ENV_FILE" ] || fail '.env is missing'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 | \
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest found'

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

for compose in "${manifests[@]}"; do
  docker compose --profile '*' --env-file "$ENV_FILE" -f "$compose" config --format json \
    | python3 -c 'import json,sys
x=json.load(sys.stdin)
for service,spec in x.get("services",{}).items():
    for p in spec.get("ports",[]) or []:
        if not isinstance(p,dict) or p.get("published") is None:
            continue
        host_ip=p.get("host_ip") or "0.0.0.0"
        published=str(p["published"])
        protocol=str(p.get("protocol","tcp"))
        print(f"{host_ip}\t{published}\t{protocol}\t{service}\t{sys.argv[1]}")' "$compose" >> "$tmp"
done

python3 - "$tmp" <<'PY'
import ipaddress
import sys


def interval(value):
    s=str(value).strip()
    if '-' in s:
        a,b=s.split('-',1)
        return int(a),int(b)
    n=int(s)
    return n,n


def wildcard(ip):
    return ip in {'0.0.0.0','::','*'}

rows=[]
with open(sys.argv[1],encoding='utf-8') as fh:
    for line in fh:
        ip,published,proto,service,compose=line.rstrip('\n').split('\t',4)
        try:
            lo,hi=interval(published)
        except ValueError as exc:
            raise SystemExit(f'invalid published port range {published!r} in {compose}: {exc}')
        if not (1 <= lo <= hi <= 65535):
            raise SystemExit(f'invalid published port range {published!r} in {compose}')
        rows.append((ip,lo,hi,proto,service,compose,published))

for i,a in enumerate(rows):
    aip,alo,ahi,aproto,asvc,acmp,arange=a
    for b in rows[i+1:]:
        bip,blo,bhi,bproto,bsvc,bcmp,brange=b
        if aproto != bproto or ahi < blo or bhi < alo:
            continue
        if not (aip == bip or wildcard(aip) or wildcard(bip)):
            continue
        raise SystemExit(
            'duplicate production host binding: '
            f'{aip}:{arange}/{aproto} ({asvc} in {acmp}) conflicts with '
            f'{bip}:{brange}/{bproto} ({bsvc} in {bcmp})'
        )

print('Production host-port ownership: PASS')
PY

log 'Production host-port ownership: PASS'
