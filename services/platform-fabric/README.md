# FTN Platform Fabric

Production foundation shared by every FTN service. Existing service implementations remain intact; this layer defines the common interoperability contract.

## Common contract

Every service should expose:

- `GET /healthz` for liveness
- `GET /readyz` for dependency-aware readiness
- `GET /metrics` for operational metrics
- authenticated service identity
- resource-scoped authorization
- correlation/request IDs
- bounded timeouts and graceful shutdown
- idempotency for retryable mutations
- audit events for state-changing operations
- explicit dependency health and version compatibility

## Security plane

- device/service identity
- certificate lifecycle hooks
- mTLS policy boundary
- audit events
- access-policy enforcement
- approval gate for privileged changes

## Platform fabric

- REST API contract boundary
- WebSocket/event contract boundary
- service/device metrics contract
- health/readiness signals
- correlation IDs
- centralized service registration

## Compatibility rule

This foundation is additive. Existing endpoints, data, services, and configuration are preserved. New integrations must use versioned contracts rather than silently changing an existing contract.

Private keys and production credentials never belong in Git.
