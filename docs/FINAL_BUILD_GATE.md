# FTN Final Build Gate

This repository uses an explicit implementation gate instead of treating architecture/configuration as a production release.

## Required components

- Go control-plane module and dependencies
- PostgreSQL schema/migrations
- Protected FTN API surface
- Web frontend
- Android client/build pipeline
- Router/core-router agent integration
- MikroTik/OLT adapters where supported
- DNS/traffic/latency collectors
- AI/Encoder/Codec service boundaries
- Integration and smoke tests
- One-click server bootstrap

## Release rule

A release may be marked production-ready only after a clean supported host can install dependencies, validate Compose, apply the database schema, build the control plane, start required services, pass health checks, and complete the applicable integration tests.

Provider-specific integrations must use documented/public interfaces and credentials supplied outside Git. New services are extension modules rather than hard-coded assumptions in the bootstrap.

## Verification commands

```bash
# Go service
cd services/control-plane
go mod download
go test ./...
go build ./...

# Compose
cd ../..
docker compose -f services/control-plane/docker-compose.yml config
docker compose -f services/control-plane/docker-compose.yml build
docker compose -f services/control-plane/docker-compose.yml up -d
```

If any command fails, fix the underlying implementation and repeat the gate. Do not label the project production-complete based only on configuration files.
