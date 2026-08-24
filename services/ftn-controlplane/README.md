# FTN Control Plane

Production service boundary for FTN identity, service registry, entitlement, device requests and durable control-plane work.

## Contracts

- `/healthz` — liveness
- `/readyz` — readiness
- `/api/v1/services` — service catalog
- `/api/v1/entitlements` — entitlement projection
- `/api/v1/device-requests` — authorized device-service request boundary

Privileged device operations remain approval-gated and auditable. Firmware deployment is optional; firmware artifacts can be built and signed without automatic deployment.
