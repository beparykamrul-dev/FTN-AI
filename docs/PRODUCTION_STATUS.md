# FTN Production Status

## Architecture/core contract

Complete for the currently defined FTN Main boundaries: Private AI, extensible AI layers, Encoder, Codec, identity/permissions, billing/payment fabric, control-plane integration, audit, and deployment extension points.

## Runnable production implementation

**Not yet certified complete.** Certification requires the Go backend, database migrations, frontend/Android clients, service integrations, collectors, and integration tests to build and run successfully on a clean supported server.

## One-click deployment

`deploy/one-click/bootstrap.sh` prepares a supported Debian/Ubuntu-style host, validates Docker Compose configuration, creates runtime directories, and starts the stack when a Compose manifest is present.

## Production gate

Before declaring FTN live, verify:

- `docker compose config` succeeds
- all required images/builds are available
- database migrations apply cleanly
- backend health checks pass
- frontend and Android builds succeed
- payment webhooks are authenticated and idempotent
- TLS/DNS/firewall are configured
- secrets are supplied outside Git
- backups and restore tests pass
- monitoring and audit pipelines are receiving events
- integration tests pass on a clean environment

No feature or service is considered production-complete merely because its architecture/configuration has been committed.
