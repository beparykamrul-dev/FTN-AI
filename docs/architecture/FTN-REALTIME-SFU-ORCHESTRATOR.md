# FTN Realtime SFU Orchestrator

## Scope
FTN realtime media fabric integrating LiveKit, mediasoup and coturn behind a provider-neutral orchestration boundary.

## Roles
- LiveKit: primary Go/WebRTC SFU and realtime AI/media path.
- mediasoup: optional secondary SFU adapter for workloads requiring its lower-level media controls.
- coturn: STUN/TURN NAT traversal and relay fallback.

## Orchestration
1. Client connects to FTN Realtime Gateway.
2. Health and Geo/S2 policy scores candidate realtime POPs.
3. Orchestrator selects a healthy SFU/TURN path.
4. Session telemetry records RTT, jitter, packet loss, bitrate and connection state.
5. Degradation triggers adaptive subscription/quality policy.
6. Capacity thresholds trigger horizontal scale-out.
7. Node failure triggers session migration/recovery where supported.
8. Recovery is verified before the node returns to service.

## Data and observability
- OpenTelemetry for traces/telemetry.
- Prometheus-compatible metrics.
- ClickHouse for high-volume realtime analytics.
- Redis only for distributed coordination/session state where required.
- Event-driven health and capacity signals feed the FTN AIOps layer.

## Security
- TLS/mTLS for control and service-to-service paths.
- Short-lived authenticated session credentials.
- Provider isolation: no provider-specific credential is embedded in the core fabric.

## Design principle
The FTN core remains provider-neutral. LiveKit, mediasoup and coturn are adapters/capabilities selected by policy, health, latency and workload rather than hard-coded into every service.
