# FTN Control Plane

Production service boundary for FTN identity, service registry, entitlement, device requests and durable control-plane work.

## Contracts

- `/healthz` — liveness
- `/readyz` — readiness
- `/api/v1/services` — service catalog
- `/api/v1/entitlements` — entitlement projection
- `/api/v1/device-requests` — authorized device-service request boundary
- `/api/v1/events/append` — authenticated durable event append
- `/api/v1/events` — tenant-scoped event replay
- `/api/v1/events/offset` — consumer offset read
- `/api/v1/events/offset/commit` — monotonic consumer offset commit

`EventJournal` remains the stable service contract. `MemoryEventJournal` is retained for deterministic tests/local development, while `RemoteEventJournal` provides the production adapter to the FTN control-plane durable journal without coupling this service to a PostgreSQL driver.

Privileged device and network operations remain approval-gated and auditable. Durable execution must preserve correlation IDs, tenant scope, monotonic offsets, idempotency and retry-safe semantics.
