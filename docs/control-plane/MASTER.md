# FTN Universal Control Plane — Master

The Control Plane unifies FTN services, infrastructure, DNS, databases, certificates, mesh, consistency, observability, visualization, AI, application delivery, and controlled operations.

## Service platform

- FTN AI Agent
- FTN DNS
- FTN Metrics
- FTN API
- FTN Proxy
- FTN WebSocket / Socket
- FTN Service Runtime / Registry
- FTN Observability / Incident / Diagnostic
- FTNMesh / FTNConsistency
- AI Call Center
- AI Billing Gateway
- FTN Messenger
- FTN Mail / Notification / File / Search / Workflow
- Dashboard Builder / Android Builder / App Builder
- App Registry / FTN App Store

## Infrastructure operations

- Unified server fabric
- Capacity, health, latency and failure-domain awareness
- DNS resolver/domain monitoring
- Database monitoring
- ACME certificate lifecycle monitoring
- GIS and topology visualization
- Service/dependency graphs
- Incident and root-cause drill-down
- Source catalog and on-demand deployment

## Operational lifecycle

```text
Source → Validate → Build → Analyze → Select Target → Policy → Approval → Deploy → Verify → Register → Audit
```

## Dashboard lifecycle

```text
Global Summary → Module → Service → Node → Interface/MAC/IP → Metrics → Logs/Trace → Dependency → Incident → Root Cause → AI Advisory
```

## Security invariants

- FTN Metrics/API first
- Deterministic logic first
- AI is advisory only
- Privileged actions are backend policy controlled
- Approval is required for privileged operations
- No secrets in telemetry or UI
- No external telemetry export
- No destructive migrations
- No blind restart/failover
- Immutable audit trail

## GIS boundary

Map visualization may use Mapbox for geographic presentation, but FTN private telemetry, secrets, and credentials are never exported to the GIS provider.

## UI principles

The Control Center is responsive, accessible, real-time, progressive-disclosure based, and designed for high signal-to-noise operational work. Overview pages summarize; module pages provide detailed monitoring and approved actions.
