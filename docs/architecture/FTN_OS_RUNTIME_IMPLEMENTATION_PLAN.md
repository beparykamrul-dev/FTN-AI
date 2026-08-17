# FTN.OS Runtime Implementation Plan

## Scope
Turn the durable-runtime contract into production Go services without claiming completion until the implementation and validation gates pass.

## Runtime components

1. `job` — durable job lifecycle, idempotency, retries, checkpoints, cancellation and DLQ.
2. `coordination` — lease locks, fencing tokens, leader election and split-brain protection.
3. `events` — append-only journal, sequence numbers, replay and consumer offsets.
4. `reconcile` — desired/observed state comparison and idempotent convergence.
5. `audit` — immutable action records and correlation across jobs/events/traces.

## Required interfaces

- `JobStore`
- `LeaseStore`
- `EventStore`
- `StateStore`
- `AuditStore`
- `Clock`
- `IDGenerator`

## Production invariants

- Every mutating job has an idempotency key.
- Every distributed lock has a lease and fencing token.
- A stale leader cannot commit after lease loss.
- Event consumers persist offsets transactionally with their state where supported.
- Reconciliation actions are idempotent.
- Cross-tenant resources are never accessible through a tenant-scoped operation.
- Destructive operations require explicit approval according to policy.
- Retries are bounded and observable.
- Restart must resume from the latest durable checkpoint.

## Test gates

- Unit tests for lifecycle/state transitions.
- Duplicate-delivery/idempotency tests.
- Lease expiry and fencing tests.
- Leader failover tests.
- Event replay tests.
- Restart/checkpoint recovery tests.
- Network-partition tests.
- Cross-tenant authorization tests.
- Load and concurrency tests.
- Security scan before production release.

## Status

This document is an implementation plan. It does not mark the runtime as production-ready.
