# FTN Production Release Gate

This gate is intentionally fail-closed. FTN must not be labelled production-complete until every required verification passes.

## Backend

- [ ] Clean-environment Go build succeeds.
- [ ] Required `go.mod` / module dependencies resolve.
- [ ] Auth middleware is enforced on protected routes.
- [ ] RBAC is enforced server-side.
- [ ] Database migrations apply cleanly.
- [ ] Accounts, billing, payments and NOC APIs are integrated with the running control plane.
- [ ] Health and readiness endpoints pass.
- [ ] Integration tests pass.

## Control Panel

- [ ] Runtime entrypoint builds successfully.
- [ ] Authentication and route guards are active.
- [ ] Dashboard receives live backend data.
- [ ] Accounts, billing, payments and NOC views are wired to APIs.
- [ ] WebSocket authentication and authorization are verified.

## Network

- [ ] Device enrollment is authenticated and authorized.
- [ ] MikroTik adapter is integration-tested against an authorized device.
- [ ] OLT/ONU drivers are integration-tested.
- [ ] SNMP/traffic telemetry is collected and normalized.
- [ ] Discovery and topology workflows pass integration tests.
- [ ] Recovery actions require explicit approval and produce an audit event.

## Clients

- [ ] Android API contract is validated against the running backend.
- [ ] Android release build succeeds.
- [ ] TV/PC clients, where enabled, build successfully.

## Security

- [ ] Secrets are externalized; credentials are never sent to clients or WebSocket payloads.
- [ ] TLS/mTLS requirements are verified for production endpoints.
- [ ] Audit trail is enabled for privileged operations.
- [ ] Destructive AI/network operations are approval-gated.
- [ ] Backup and restore are tested.

## Release rule

A release is **NOT production-complete** if any mandatory checkbox above remains unverified. Passing configuration or contract checks alone is insufficient; clean builds and integration tests are required.
