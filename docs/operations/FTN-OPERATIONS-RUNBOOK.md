# FTN Operations Runbook

## Purpose

This runbook is the operational baseline for FTN production. It turns the architecture/release contracts into repeatable day-2 operations without making the FTN Main server a mandatory traffic transit point.

## Golden operating loop

`Observe -> Assess -> Act (policy/approval) -> Verify -> Audit`

## Daily checks

- Control-plane health
- POP health and capacity
- DNS authoritative/recursive health
- BGP session and route-policy health
- Mesh convergence and stale-peer state
- Gateway/VPN health
- Database replication and backup freshness
- Telemetry ingestion/retention
- Certificate/PKI expiry and revocation state
- Security events and unresolved incidents
- Customer SLO/SLA degradation

## Deployment rules

1. Validate configuration and contracts before rollout.
2. Deploy to one canary node/POP first.
3. Drain affected services before disruptive changes.
4. Observe health, RTT, jitter, loss and error rate.
5. Promote only after verification.
6. Roll back on defined failure thresholds.
7. Record the change and observed result.

## Failure handling

### POP failure

`Detect -> isolate -> select healthy POP/path -> restore critical services -> verify -> incident record`

### Provider failure

`Detect -> adapter health check -> alternate eligible provider/path -> verify -> cost/policy check -> audit`

### Database failure

`Detect -> preserve writes/queue where safe -> fail over to healthy replica -> integrity check -> resume -> audit`

### Certificate/PKI failure

`Detect -> deny unsafe service access -> use approved rotation/revocation workflow -> verify identity -> restore`

### Security incident

`Detect -> correlate -> contain according to authorized policy -> preserve evidence -> verify containment -> escalate/report -> recover`

## Rollback requirements

Every production change must have a known rollback or recovery path. Database migrations require tested backup/restore or reversible migration strategy before promotion.

## Capacity rules

Workloads should prefer the best eligible node based on workload compatibility, health, locality, latency, capacity, quota and policy. Free capacity is not sufficient by itself to move data or services.

## Mesh rules

The mesh is for distributed service/traffic sharing. Peer connectivity does not imply unrestricted trust. Identity, authorization, segmentation and service-level policy remain mandatory.

## Security rules

- Default deny for privileged service access.
- Least privilege and tenant isolation.
- mTLS for trusted service-to-service communication.
- Audit privileged mutations.
- No unauthorized surveillance or access.
- No autonomous destructive action without an explicitly authorized policy.

## Production evidence

A deployment is not considered production-verified from source code alone. Real authorized routers, servers, DNS endpoints, BGP peers, PKI services and WAN paths must pass acceptance tests with measured results.

## Change ownership

The registry and versioned contracts are the source of truth. Provider-specific implementations remain behind adapters. New services should be additive and must preserve policy, security, observability, verification and audit boundaries.
