# FTN Control Plane

Production service boundary for the FTN service catalog, entitlement-aware clients, device/service requests, and PostgreSQL persistence.

## Runtime

- `GET /healthz` — liveness
- `GET /readyz` — readiness and database check
- `GET /api/v1/services` — service catalog
- `GET /api/v1/entitlements` — entitlement projection from `X-FTN-Services`
- `POST /api/v1/service-requests` — authorized service/device request intake

Firmware deployment is deliberately disabled by default. Requests are persisted and returned as `202 Accepted`; any future device change must pass authorization, compatibility, signing, and approval gates.

## Local production-like stack

Set `FTN_DB_PASSWORD` and run:

```sh
docker compose -f services/control-plane/docker-compose.yml up --build -d
```

The stack contains PostgreSQL and the control-plane service. The frontend is served by the control-plane root route in the current lightweight deployment boundary.
