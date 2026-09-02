#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE="${ROOT}/deployment/dns/docker-compose.yml"

command -v docker >/dev/null 2>&1 || { echo "Docker is required" >&2; exit 1; }
docker compose version >/dev/null 2>&1 || { echo "Docker Compose v2 is required" >&2; exit 1; }

cd "${ROOT}"

# Core runtime: dnsdist + Unbound + CoreDNS.
docker compose -f "${COMPOSE}" config >/dev/null
docker compose -f "${COMPOSE}" up -d ftn-unbound ftn-coredns ftn-dnsdist

# Optional management/ACME components are explicit so a licensed Enterprise
# image or DNS-provider credentials can be supplied without putting secrets in git.
if [[ "${FTN_DNS_INSTALL_TECHNITIUM:-0}" == "1" ]]; then
  docker compose -f "${COMPOSE}" --profile enterprise up -d ftn-technitium
fi

if [[ "${FTN_DNS_INSTALL_CADDY:-0}" == "1" ]]; then
  docker compose -f "${COMPOSE}" --profile acme up -d ftn-caddy
fi

echo "FTN DNS core runtime started."
echo "Verify with: docker compose -f deployment/dns/docker-compose.yml ps"
echo "Do not advertise Anycast/BGP until DNS health and route convergence have been verified in the target network."
