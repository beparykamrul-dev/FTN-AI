#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE="${ROOT}/deployment/dns/docker-compose.yml"
DNSDIST_CONF="${ROOT}/configs/dns/dnsdist.conf"

command -v docker >/dev/null 2>&1 || { echo "Docker is required" >&2; exit 1; }
docker compose version >/dev/null 2>&1 || { echo "Docker Compose v2 is required" >&2; exit 1; }

cd "${ROOT}"
echo "[1/5] validating Compose"
docker compose -f "${COMPOSE}" config >/dev/null

echo "[2/5] checking core containers"
for service in ftn-unbound ftn-coredns ftn-dnsdist; do
  state="$(docker compose -f "${COMPOSE}" ps -q "${service}" | xargs -r docker inspect -f '{{.State.Status}}' 2>/dev/null || true)"
  [[ "${state}" == "running" ]] || { echo "${service} is not running" >&2; exit 1; }
done

echo "[3/5] checking Unbound health"
unbound_health="$(docker compose -f "${COMPOSE}" ps -q ftn-unbound | xargs -r docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}')"
[[ "${unbound_health}" == "healthy" ]] || { echo "Unbound health: ${unbound_health}" >&2; exit 1; }

echo "[4/5] checking DNS listener from host"
if command -v dig >/dev/null 2>&1; then
  dig +time=3 +tries=1 @127.0.0.1 familytimenet.com SOA >/dev/null
  dig +time=3 +tries=1 @127.0.0.1 familytimenet.com NS >/dev/null
else
  echo "dig not installed; container/runtime checks completed, external DNS query skipped" >&2
fi

echo "[5/5] validating configured DNS listener claims"
# A published port is not evidence that a protocol is configured. Require an
# explicit listener declaration in dnsdist.conf before testing optional ports.
required_tcp_ports=(53)
for port in "${required_tcp_ports[@]}"; do
  grep -Eq "(^|[^0-9])${port}([^0-9]|$)" "${DNSDIST_CONF}" || {
    echo "dnsdist TCP/UDP listener for port ${port} is not declared" >&2
    exit 1
  }
done

for optional in 853 443 784; do
  case "${optional}" in
    853) pattern='listen.*:853|listen.*853' ;;
    443) pattern='listen.*:443|listen.*443' ;;
    784) pattern='listen.*:784|listen.*784' ;;
  esac
  if grep -Eiq "${pattern}" "${DNSDIST_CONF}"; then
    echo "configured optional listener detected: ${optional}"
  else
    echo "optional listener ${optional}: not configured (not treated as live)"
  fi
done

echo "FTN DNS runtime validation passed for configured listeners."
echo "Public delegation, Anycast/BGP, provider failover, secure DNS transports, and multi-site recovery require target-network acceptance tests."
