#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

contract="configs/v1/ftn-service-native-multiprotocol.yaml"
registry="configs/transport/ftn-protocol-registry.yaml"
[[ -s "$contract" ]] || { echo "missing native multi-protocol contract" >&2; exit 1; }
[[ -s "$registry" ]] || { echo "missing protocol registry" >&2; exit 1; }

for service in control api web dns ddns proxy socket websocket mesh buckboon routed metrics monitoring telemetry silk ai billing payment iptv support notification device_driver; do
  grep -Eq "^    ${service}: \[" "$contract" || { echo "missing protocol matrix for ${service}" >&2; exit 1; }
done

grep -q 'protocol_identity: ftn-native' "$contract"
grep -q 'services_are_protocol_endpoints: true' "$contract"
grep -q 'capability_advertisement: true' "$contract"
grep -q 'health_selection: true' "$contract"
grep -q 'security_downgrade: prohibited' "$contract"
grep -q '    primary: silk' "$contract"
grep -q 'default_action: deny-unregistered' "$registry"

echo 'FTN native multi-protocol gate: PASS'
