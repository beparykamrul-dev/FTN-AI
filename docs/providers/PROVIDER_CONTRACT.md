# FTN Provider Contract

## Purpose

FTN integrates DNS, network, observability, security, and infrastructure providers through a stable adapter boundary. External projects and services are treated as capabilities to evaluate, not code to copy.

## Provider requirements

Every provider adapter should define:

- capability and protocol identity
- supported read operations
- normalized resource model
- health and latency probes
- timeout and cancellation behavior
- authentication boundary
- retry/idempotency expectations
- error classification
- version/compatibility information
- audit requirements for mutations

## Read path

```text
Provider
  -> Adapter
  -> Normalizer
  -> Snapshot
  -> Consistency Engine
  -> Control Plane
```

## Mutation path

```text
Plan -> Audit -> Explicit Approval -> Guard -> Idempotency -> Executor
```

Provider adapters must not bypass the approval boundary.

## External technology evaluation

Projects such as POLER or AOS.NET may be added only after their exact upstream identity, license, protocol, security model, and operational behavior have been verified. Unknown or ambiguous names remain research-only entries and are not production dependencies.
