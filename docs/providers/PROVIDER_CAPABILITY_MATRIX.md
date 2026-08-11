# FTN Provider Capability Matrix

The capability matrix is the machine-readable admission model for provider integrations.

## Capability dimensions

| Capability | Meaning |
|---|---|
| import | Read and normalize authoritative resources |
| snapshot | Produce deterministic state snapshots |
| health | Report provider health |
| latency | Measure endpoint latency |
| reconcile | Participate in approved state reconciliation |
| dnssec | DNSSEC-aware operations |
| anycast | Anycast/POP-aware routing metadata |
| api | Native API integration |
| audit | Mutation audit support |

## Admission rules

A provider must explicitly declare supported capabilities. Unknown capabilities are treated as unsupported. Read-only capabilities may be used by the consistency and monitoring layers without granting mutation authority.

Mutation-related capabilities require the complete FTN execution boundary:

```text
Plan -> Audit -> Approval -> Guard -> Idempotency -> Executor
```

## Versioning

Capability changes must be versioned. Removing a capability is a compatibility event and must follow the FTN deprecation policy.
