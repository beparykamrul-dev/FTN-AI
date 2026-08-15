# FTN Infrastructure Platform

## Architecture decision
FTN may use a compatible open-source virtualization foundation where licensing permits, while keeping FTN-specific orchestration, network automation, service catalog, metrics, topology, AI policy and user experience in FTN-owned layers.

## Layer model

```text
Hardware
  -> Virtualization / Container Foundation
  -> FTN Node Agent
  -> FTN Control Plane
  -> FTN Network / FTNWAN
  -> FTN API + Metrics
  -> Discovery + Topology
  -> Fiber Map
  -> Portals / Apps / AI Assistant
```

## FTN-owned control plane
- Node registration and health
- Cluster inventory
- Resource balancing policy
- Network intent and configuration state
- Device discovery
- Metrics and audit events
- Service access policy
- AI recommendations and approval workflow

## Safety boundary
AI recommendations must pass authorization, policy and audit checks before any network or infrastructure mutation. No AI component should receive unrestricted administrator credentials by default.

## Open-source boundary
Keep upstream open-source components identifiable and license-compliant. FTN customization should be implemented as integrations, plugins, agents, APIs and independent control-plane components where practical rather than obscuring upstream provenance.

## Production objective
The result should feel like one FTN platform to operators while retaining a clear technical and legal boundary between upstream open-source foundations and FTN-owned components.
