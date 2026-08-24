# FTN Dynamic Scale Policy

## Scope

FTN infrastructure is designed without a hard-coded server-count or storage-capacity ceiling in the application architecture.

The Control Plane treats Main, POP, Client, Backup, DNS, Edge, Hosting, Cache, Storage, and auxiliary nodes as dynamically discovered resources.

## Scaling model

```text
Node Enrollment
    -> Identity / Authorization
    -> Capability Discovery
    -> Health Check
    -> Capacity Observation
    -> Placement / Routing Policy
    -> Reconciliation
    -> Monitoring
```

Capacity is observed at runtime. The Control Plane must not encode a fixed maximum number of servers or a fixed storage-pool size.

## Node lifecycle

- Discover authorized nodes.
- Register node identity and capabilities.
- Continuously observe health and available resources.
- Add eligible nodes without rebuilding the core architecture.
- Rebalance services according to policy and health.
- Drain and remove nodes safely.
- Preserve audit and recovery information.

## Service scaling

Services use capability-based placement rather than fixed host lists. A service may run on one or many eligible nodes depending on its deployment policy, health, locality, and available capacity.

## Control Panel

The Control Panel exposes current infrastructure state, including:

- node count
- node roles
- health
- CPU and memory availability
- storage capacity and utilization
- network capacity
- service placement
- replication state
- failover state

These are runtime observations, not architectural limits.

## Safety

Automatic scaling and reconciliation remain subject to authorization, policy, health verification, and audit logging. The absence of a hard-coded capacity limit does not authorize arbitrary infrastructure or network changes.
