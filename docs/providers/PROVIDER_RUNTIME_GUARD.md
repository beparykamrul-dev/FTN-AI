# FTN Provider Runtime Guard

Registry admission is a build-time and governance decision; runtime authorization is a separate control.

## Runtime flow

```text
Incoming Operation
      ↓
Provider Registry Lookup
      ↓
Capability Check
      ↓
Authentication / Scope
      ↓
Plan Validation
      ↓
Audit
      ↓
Explicit Approval
      ↓
Idempotency Key
      ↓
Runtime Guard
      ↓
Provider Executor
```

## Guard requirements

- reject unknown providers
- reject unsupported capabilities
- enforce operation scope
- require an approved plan for mutations
- require a valid idempotency key
- enforce timeout and cancellation
- preserve audit correlation
- fail closed on missing security metadata

Read-only health, latency, and import operations may bypass mutation approval only when their registered capability explicitly permits read access.

The guard is the final policy boundary immediately before a provider executor and must not be implemented inside individual provider adapters.
