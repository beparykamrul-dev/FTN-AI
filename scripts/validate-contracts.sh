#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

required=(
  "docs/PROJECT_COMPLETION_AZ.md"
  "docs/DEPLOYMENT_ACCEPTANCE.md"
  "configs/v1/ftn-connectivity-stack.yaml"
  "configs/v1/ftn-kernel-backend.yaml"
  "configs/v1/ftn-gateway-fabric.yaml"
  "configs/v1/ftn-global-latency.yaml"
  "configs/v1/ftn-global-dns-mesh.yaml"
  "internal/platform/router/router.go"
  "internal/platform/dns/familytimenet_node_registry.go"
  "frontend/control-center/index.html"
  "frontend/control-center/components/component-contracts.yaml"
)

for path in "${required[@]}"; do
  [[ -f "$path" ]] || { echo "missing required contract: $path" >&2; exit 1; }
done

# Prevent accidental committed credentials/private keys in source-level CI.
if git grep -n -I -E '(BEGIN (RSA|EC|OPENSSH|PRIVATE) KEY|AKIA[0-9A-Z]{16}|ghp_[A-Za-z0-9]{20,}|sk-[A-Za-z0-9_-]{20,})' -- ':!docs/DEPLOYMENT_ACCEPTANCE.md' ':!docs/PROJECT_COMPLETION_AZ.md'; then
  echo "possible credential/private-key material found" >&2
  exit 1
fi

# The canonical DNS contract must remain stable and registry-driven.
grep -q 'familytimenet.com' configs/v1/ftn-global-dns-mesh.yaml
grep -q 'registryDriven: true' configs/v1/ftn-global-dns-mesh.yaml
grep -q 'additiveNodes: true' configs/v1/ftn-global-dns-mesh.yaml

# The control-center must remain service-catalog driven rather than hard-coded.
grep -q '/api/v1/services' frontend/control-center/index.html
grep -q 'FTN Control Center' frontend/control-center/index.html

echo 'FTN-AI contract validation passed.'
