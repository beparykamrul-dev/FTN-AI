# FTN Implementation Audit

## Current verified state

The repository contains an executable-looking Go control-plane source with HTTP health/readiness endpoints, a service catalog, service-request persistence, node/network/ex-router routes, security headers, and an embedded HTML frontend.

## Not yet verified as production-complete

The current repository state does **not** provide enough evidence to claim that the complete backend/frontend/router stack has been build-tested successfully.

### Backend gaps requiring validation

- A root `go.mod` was not found during repository search, while `services/control-plane/main.go` imports `github.com/jackc/pgx/v5/pgxpool`.
- The referenced helper implementations (`metricsHandler`, `a.nodeCatalog`, `a.registerNode`, `a.nodeHeartbeat`, `a.placement`, `a.networkHealth`, `a.exRouter`) must exist and compile together.
- Billing/account/permission APIs need to be wired into the control-plane code rather than existing only as configuration/database contracts.
- Authentication/authorization middleware must be enforced on protected routes.

### Frontend gaps requiring validation

- The current frontend is an embedded HTML service catalog, not evidence of a complete Web application.
- Android/TV/PC client builds are not verified here.

### Router/core-router gaps requiring validation

- No verified core-router/local-router firmware or router-agent implementation was found by repository code search.
- MikroTik/OLT integration, secure device enrollment, firmware build/signing, and authorized push/rollback require actual implementation and integration tests.

## Rule

Do not mark FTN production-complete until a clean environment can build the backend and clients, apply migrations, start the services, pass health/readiness and integration tests, and validate router/device workflows.
