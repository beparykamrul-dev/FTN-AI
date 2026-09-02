#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE="${ROOT}/deployment/dns/docker-compose.yml"

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

echo "[5/5] checking exposed sockets"
docker compose -f "${COMPOSE}" port ftn-dnsdist 53 >/dev/null
docker compose -f "${COMPOSE}" port ftn-dnsdist 853 >/dev/null
docker compose -f "${COMPOSE}" port ftn-dnsdist 443 >/dev/null
docker compose -f "${COMPOSE}" port ftn-dnsdist 784 >/dev/null

echo "FTN DNS runtime validation passed."
echo "This validates the local container stack only; public delegation, Anycast/BGP, provider failover, and multi-site recovery require tests in the target FTN network."
