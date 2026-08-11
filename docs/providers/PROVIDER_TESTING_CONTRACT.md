# FTN Provider Testing Contract

Every provider and source-wrapper integration must pass the same contract suite before production registration.

## Required test classes

1. Capability discovery
2. Read/import normalization
3. Snapshot determinism
4. Health and latency probing
5. Timeout and cancellation
6. Authentication failure handling
7. Retry and idempotency behavior
8. Drift detection
9. Audit and approval enforcement
10. Provider-specific error classification

## Mutation safety

Tests must prove that import, normalization, health checks, and drift analysis cannot mutate authoritative provider state. Mutation tests must exercise the explicit plan -> audit -> approval -> guard -> idempotency -> executor path.

## Compatibility

Provider adapters should expose a version/capability matrix. Unsupported features must fail explicitly rather than silently degrading authoritative state.

## CI gate

A provider is eligible for the FTN registry only when its contract tests pass and its security/license metadata has been reviewed.
