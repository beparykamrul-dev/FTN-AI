# FTN-AI Deployment Acceptance

## Source-level gate

- [x] Versioned architecture/configuration contracts
- [x] Kernel/backend/router/gateway boundaries
- [x] Registry-driven DNS and service onboarding
- [x] Local/global full-mesh model
- [x] GoBGP/GeoIP2/Anycast integration boundaries
- [x] PKI/mTLS and approval boundaries
- [x] Latency/path-selection model
- [x] NetSA/ClickHouse telemetry boundary
- [x] A–Z completion contract
- [x] Repository validation workflow

## Infrastructure acceptance gate

Before declaring a live deployment production-ready, the operator must validate against authorized infrastructure:

1. DNS authoritative/recursive behavior for `familytimenet.com`.
2. GoBGP sessions and route-policy correctness.
3. Anycast advertisement and controlled withdrawal.
4. Local and global mesh convergence.
5. Gateway/proxy authentication and egress policy.
6. PKI issuance, rotation and revocation.
7. Router/device adapter connectivity.
8. End-to-end RTT, jitter, loss and p95/p99 latency under load.
9. Failure and drain scenarios.
10. Telemetry delivery and retention.

No simulated result, placeholder endpoint, embedded credential or fake health signal may be used to satisfy this gate.

## Completion definition

The repository is source-complete when CI passes and the contracts are internally consistent. The infrastructure gate is complete only after real authorized systems pass the acceptance tests above.
