# FTN Provider Registry

The provider registry is the production admission boundary for DNS, network, observability, security, and infrastructure integrations.

## Admission flow

```text
Candidate Provider
      ↓
Identity + License Review
      ↓
Capability Contract
      ↓
Contract Tests
      ↓
Security Review
      ↓
Compatibility Check
      ↓
Registry Admission
```

## Registry metadata

Each admitted provider should record:

- canonical provider name and upstream identity
- adapter version
- protocol/capability set
- supported read and mutation operations
- authentication requirements
- health/latency probe support
- compatibility range
- security/license review status
- operational owner

## Safety rules

Unknown or unverified integrations remain research-only. A provider cannot become executable merely by appearing in configuration. Mutation capabilities must still pass the FTN plan, audit, explicit approval, guard, and idempotency boundaries.
