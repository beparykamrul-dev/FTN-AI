#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ROOT_DIR="${FTN_COMPOSE_ROOT:-$SCRIPT_ROOT}"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"
source "$SCRIPT_ROOT/scripts/production-compose-env.sh"

fail(){ printf '[FTN][HOST-PORTS][ERROR] %s\n' "$*" >&2; exit 1; }
log(){ printf '[FTN][HOST-PORTS] %s\n' "$*"; }
trap 'printf "[FTN][HOST-PORTS][ERROR] failed at line %s\n" "$LINENO" >&2' ERR

command -v docker >/dev/null 2>&1 || fail 'Docker is required'
command -v python3 >/dev/null 2>&1 || fail 'python3 is required'
command -v ss >/dev/null 2>&1 || fail 'iproute2/ss is required'
[ -d "$ROOT_DIR" ] || fail "Compose root does not exist: $ROOT_DIR"
[ -f "$ENV_FILE" ] || fail '.env is missing'
docker info >/dev/null 2>&1 || fail 'Docker daemon is not available'

mapfile -t manifests < <(find "$ROOT_DIR" -type f \( -name 'docker-compose.yml' -o -name 'compose.yml' \) \
  -not -path '*/.git/*' -not -path '*/node_modules/*' -print0 |
  xargs -0 -r grep -El 'FTN_PRODUCTION_STACK=true|x-ftn-production-stack:[[:space:]]*true' | sort)
((${#manifests[@]})) || fail 'No production Compose manifest found'

tmp="$(mktemp)"
ss_tmp="$(mktemp)"
trap 'rm -f "$tmp" "$ss_tmp"' EXIT

for compose in "${manifests[@]}"; do
  docker compose --env-file "$ENV_FILE" -f "$compose" config --format json |
    python3 -c 'import json,sys
x=json.load(sys.stdin)
for service,spec in x.get("services",{}).items():
    for p in spec.get("ports",[]) or []:
        if not isinstance(p,dict) or p.get("published") is None: continue
        ip=p.get("host_ip") or "0.0.0.0"; pub=str(p["published"]); proto=str(p.get("protocol","tcp")).lower()
        if "-" in pub: lo,hi=map(int,pub.split("-",1))
        else: lo=hi=int(pub)
        print(f"{ip}\t{lo}\t{hi}\t{proto}\t{service}\t{sys.argv[1]}")' "$compose" >> "$tmp"
done

ss -H -ltnup 2>/dev/null > "$ss_tmp" || true
python3 - "$tmp" "$ss_tmp" <<'PY'
import re,sys
from ipaddress import ip_address
expected=[]
with open(sys.argv[1],encoding='utf-8') as f:
    for line in f:
        ip,lo,hi,proto,service,compose=line.rstrip('\n').split('\t')
        expected.append((ip,int(lo),int(hi),proto,service,compose))

def parse_local(value):
    value=value.strip()
    if value.startswith('['):
        m=re.match(r'^\[([^]]+)\]:(\d+)$',value)
        return (m.group(1),int(m.group(2))) if m else (value,None)
    host,sep,port=value.rpartition(':')
    return (host,int(port)) if sep and port.isdigit() else (value,None)

def ip_overlap(a,b):
    if a in ('*','0.0.0.0','[::]','::') or b in ('*','0.0.0.0','[::]','::'): return True
    try: return ip_address(a)==ip_address(b)
    except ValueError: return a==b

def docker_owned(line):
    return 'users:(("docker-proxy"' in line or 'users:(("dockerd"' in line)

with open(sys.argv[2],encoding='utf-8',errors='replace') as f:
    for raw in f:
        line=raw.strip()
        if not line or docker_owned(line): continue
        fields=line.split()
        if len(fields)<5: continue
        proto=fields[0].lower(); host,port=parse_local(fields[4])
        if port is None: continue
        for eip,elo,ehi,eproto,service,compose in expected:
            if proto != eproto or not (elo <= port <= ehi): continue
            if ip_overlap(host,eip):
                raise SystemExit(f'host listener conflict: {host}:{port}/{proto} is used by a non-Docker process; expected by {service} in {compose}')
print('Host process port ownership: PASS')
PY

log "Checked Compose root: $ROOT_DIR"
log 'Host process port ownership: PASS'
