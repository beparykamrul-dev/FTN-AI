# FTN-AI Production Completion

## Status

FTN-AI is organized as a modular, provider-neutral platform. The repository README defines the AI layer as private, policy-controlled, approval-first, and intended for real FTN infrastructure rather than demo-only implementations.

## Permanent architecture boundaries

- Core control plane: identity, PKI, policy, service registry, configuration and orchestration.
- Network plane: DNS, routing, mesh, traffic intelligence and edge connectivity.
- Security plane: Zeek/FastNetMon/SIEM/eBPF/XDP adapters, evidence and incident workflows.
- Observability plane: OpenTelemetry, metrics, logs, flow telemetry, ClickHouse/search and NOC dashboards.
- Communication plane: FTN notifications, chat, audio/video, groups, realtime gateway and TURN/SFU adapters.
- Game plane: matchmaking, sessions, authoritative servers, QoS and game telemetry.
- Edge/cloud plane: POP services, user/developer opt-in cloud nodes, workload placement and autoscaling.
- Provider plane: replaceable adapters for external cloud/CDN/DNS/network/service providers.
- Data plane: fit-for-purpose databases and storage selected by policy, capacity, quota, latency and data requirements.
- Governance plane: RBAC, audit trail, tenant isolation, approval gates, license/compliance checks and retention policies.

## Operating model

`Control centralized -> traffic distributed -> services decentralized -> security enforced -> telemetry federated -> recovery automated`

The FTN main server must not be the mandatory transit point for ordinary mesh traffic. POPs and authorized user/developer nodes may provide local services and workload capacity subject to explicit policy, resource quotas and tenant isolation.

## Database selection

A database broker may select an appropriate local/POP/cloud backend based on workload type, available capacity, quota, latency, availability, replication state and policy. Free capacity alone must never cause an unsafe or unapproved migration.

## Notification and communication

FTN Notify/Communication is the native primary channel. Telegram, WhatsApp, IMO and other external channels remain optional interoperability/fallback adapters rather than core dependencies.

## Security and evidence

Security telemetry and digital-evidence workflows are isolated from the service core. Evidence handling requires source identity, normalized timestamps, integrity hashes, access control, retention policy and an auditable chain of custody. The platform does not autonomously declare guilt or perform unauthorized surveillance.

## AI/AIOps control loop

`Observe -> Analyze -> Recommend -> Policy/Approval Gate -> Execute -> Verify -> Audit`

Privileged, destructive or externally consequential actions remain approval-controlled unless an explicitly authorized policy permits automation.

## Production gate

Before any component is treated as production-ready, validate:

1. license and usage rights;
2. maintenance/security posture;
3. resource and dependency cost;
4. API/contract compatibility;
5. tenant and failure-domain isolation;
6. observability and health checks;
7. backup/restore behavior;
8. upgrade/rollback path;
9. security policy and auditability;
10. real hardware/network integration tests.

This document is the completion baseline; future work should implement and verify missing contracts/tests/deployments rather than continually expanding the dependency list.
