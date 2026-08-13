# FTN Mesh — Meshtastic Adapter

Meshtastic is integrated as an optional edge/low-bandwidth mesh transport adapter under FTN Mesh.

## Boundary

Meshtastic is not treated as the FTN core mesh. FTN Mesh remains the policy, identity, routing and observability layer; the adapter provides a transport boundary for compatible Meshtastic nodes.

## Integration points

- node discovery and health
- message transport
- link quality metrics
- route/path observations
- gateway status
- offline/queued delivery state
- FTN identity and authorization boundary
- FTN Metrics integration

## Safety

The adapter must not expose secrets through telemetry or export external telemetry by default. Privileged configuration changes remain policy-controlled and auditable.
