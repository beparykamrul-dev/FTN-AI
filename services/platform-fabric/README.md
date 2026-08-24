# FTN Platform Fabric

Production boundary for FTN identity/security and API/event/metrics fabric.

## Security plane
- device/service identity
- certificate lifecycle hooks
- mTLS policy boundary
- audit events
- access-policy enforcement

## Platform fabric
- REST API contract boundary
- WebSocket/event contract boundary
- service/device metrics contract
- health/readiness signals
- correlation IDs

This package contains contracts and policy boundaries only; private keys and production credentials never belong in Git.
