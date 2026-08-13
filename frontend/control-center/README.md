# FTN Universal Control Center

Production frontend foundation for the FTN Control Plane.

## UI contract

- Responsive desktop/tablet/mobile shell
- Overview summary with progressive drill-down
- Service, DNS, database, certificate, mesh and incident modules
- FTN Metrics/API-first data access
- FTN WebSocket real-time updates
- AI advisory only where analysis adds value
- RBAC and audit-aware actions
- Secrets never rendered

## Runtime data flow

```text
FTN API + FTN Metrics
        ↓
 Control Center Data Layer
        ↓
 Pages / Modules
        ↓
 FTN WebSocket events
        ↓
 Live status updates
```

## Implementation boundary

The UI is a control and observability surface. It does not bypass backend policy, execute privileged infrastructure operations directly, or expose credentials.
