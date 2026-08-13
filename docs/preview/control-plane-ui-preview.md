# FTN Universal Control Center — UI Preview

## Design direction

A polished operational console: calm visual hierarchy, compact information density, clear semantic status, responsive layouts, and fast drill-down from global summary to evidence.

## Screen map

```text
┌─────────────────────────────────────────────────────────────────────┐
│ FTN Universal Control Center       Search      Alerts   Account     │
├───────────────┬─────────────────────────────────────────────────────┤
│ Overview      │ GLOBAL HEALTH                                       │
│ Infrastructure│ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐        │
│ Services      │ │Servers │ │Service │ │ DNS    │ │Mesh    │        │
│ DNS            │ │ 128    │ │ 642    │ │99.99%  │ │98.9%   │        │
│ Database       │ └────────┘ └────────┘ └────────┘ └────────┘        │
│ Certificates   │                                                     │
│ FTNMesh        │ LATENCY / INCIDENTS / DEPLOYMENTS                  │
│ Consistency    │ ┌──────────────────────┐ ┌──────────────────────┐ │
│ Observability  │ │ Live telemetry        │ │ Active incidents     │ │
│ Incidents      │ │                      │ │                      │ │
│ AI             │ └──────────────────────┘ └──────────────────────┘ │
│ Billing        │                                                     │
│ Communication  │ SERVICE / DEPENDENCY GRAPH                         │
│ Builders       │ ┌────────────────────────────────────────────────┐ │
│ App Store      │ │ Edge → DNS → Cache → API → DB                  │ │
│ Deployments    │ └────────────────────────────────────────────────┘ │
│ Audit          │                                                     │
│ Settings       │                                                     │
└───────────────┴─────────────────────────────────────────────────────┘
```

## Interaction model

```text
Summary
  → Module
  → Service
  → Node
  → Interface / MAC / IP
  → Metric / Log / Trace
  → Dependency
  → Incident
  → Root Cause
  → AI Advisory
```

## Service detail

Every service page exposes the same operational contract:

- health and availability
- FTN Metrics
- FTN API state
- dependencies
- latency
- logs
- traces
- incidents
- deployment history
- configuration state
- audit events

## Deployment experience

```text
Source → Preview → Diff → Analyze → Target → Health Gate → Approval → Deploy → Verify
```

No blind restart, blind failover, destructive migration, or hidden privileged action is exposed in the UI.

## Responsive behavior

- Desktop: persistent navigation + multi-column operational views.
- Tablet: collapsible navigation + adaptive cards.
- Mobile: stacked cards + bottom/compact navigation + focused drill-down.
- Charts and tables resize without losing critical status information.

## Visual principles

- semantic status rather than decorative color overload
- consistent spacing and typography
- high signal-to-noise dashboards
- progressive disclosure for detailed telemetry
- keyboard/accessibility support
- light/dark themes
- real-time updates through FTN WebSocket
- no secrets rendered into the interface
