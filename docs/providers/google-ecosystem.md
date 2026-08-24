# Google Ecosystem Provider Reference

FTN provider/reference profile for Google, Google Research, Google Samples, and Google APIs.

## Capability classes
- Google APIs: REST/gRPC interface definitions, generated clients, authentication, retries, pagination, long-running operations.
- Google Research: AI/ML, data analysis, experimentation and research references.
- Google Samples: Android/iOS/web integration examples and API usage patterns.
- googleapis: API contracts, protobuf/gRPC definitions, client generation, cloud clients, authentication and API tooling.

## FTN integration boundary
Google technology is treated as a pluggable provider/reference layer, not a hard dependency of FTN Core.

Flow:

Google OSS/API/SDK -> FTN Google Adapter -> Capability Normalizer -> FTN Database -> Metrics/Telemetry -> AI Agent -> Unified Control Plane

## Relevant FTN uses
- DNS/API provider integrations
- Cloud service adapters
- gRPC/REST API gateway patterns
- OAuth2/service-account/API-key credential adapters
- AI/ML research tooling
- Android client integration references
- telemetry, retries, pagination and long-running operation patterns
- MCP/database tooling references

## Provenance and license handling
Each imported/reference component must retain repository URL, revision, license, provenance, and dependency metadata. Open-source code is only incorporated where its license permits the intended use. Provider services remain isolated behind adapters.
