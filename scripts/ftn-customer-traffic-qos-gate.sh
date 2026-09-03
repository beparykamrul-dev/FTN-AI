#!/usr/bin/env bash
set -euo pipefail

contract="configs/v1/ftn-customer-traffic-qos.yaml"
[[ -f "$contract" ]] || { echo "missing $contract" >&2; exit 1; }

for service in whatsapp telegram imo pubg freefire; do
  grep -Eq "services:.*${service}|- ${service}" "$contract" || { echo "missing service: $service" >&2; exit 1; }
done

grep -q 'availability_target: 99.99' "$contract"
grep -q 'identity_source: managed-service-registry' "$contract"
grep -q 'endpoint_ips: dynamic' "$contract"
grep -q 'engine: exrouter' "$contract"
grep -q 'health_aware: true' "$contract"
grep -q 'hysteresis: true' "$contract"
grep -q 'dscp: true' "$contract"
grep -q 'customer_fairness: true' "$contract"
grep -q 'anti_spoofing: required' "$contract"
grep -q 'policy_gated_changes: true' "$contract"
grep -q 'no_embedded_provider_credentials: true' "$contract"

echo "FTN customer traffic QoS contract: OK"
