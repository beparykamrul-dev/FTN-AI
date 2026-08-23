# FTN-AI A–Z Completion Contract

This document is the release gate for the FTN-AI repository. A feature is complete only when its contract, implementation boundary, validation, security policy, observability, and integration path are defined.

## A — Architecture
Provider-neutral control plane, kernel, backend, router, gateway and service registry.

## B — Backbone
Internet-WAN/provider mesh; no assumption of FTN-owned fiber or wireless backbone.

## C — Connectivity
WireGuard, AmneziaWG, OpenVPN 3, GRE, SSLH, Shadowsocks, Hysteria2 and WebSocket adapters.

## D — DNS
`familytimenet.com` is registry-driven. New DNS nodes/providers are additive and must implement the same contract.

## E — Edge
Anycast, reverse/edge proxy, authenticated SOCKS5 and gateway capabilities use the common gateway contract.

## F — Full Mesh
Local and global full-mesh scopes with health-aware path selection.

## G — GoBGP / GeoIP2
GoBGP is the BGP control-plane integration; GeoIP2 provides locality intelligence.

## H — Health
Active/passive probing, RTT, jitter, loss, utilization, queue depth and stability feed routing decisions.

## I — Identity
PKI/mTLS, enrollment, revocation and server-side secret boundaries.

## J — Jupyter
Kernel wrapper and ipykernel integration expose structured FTN tools, not arbitrary shell execution.

## K — Kernel
One policy-controlled kernel/tool boundary for network and service operations.

## L — Latency
End-to-end path latency is optimized, not DNS latency alone. Route changes use hysteresis to avoid flapping.

## M — Mesh
BATMAN-adv, Yggdrasil and CJDNS are optional mesh adapters behind the common topology model.

## N — NetSA
SiLK, YAF, super_mediator and PySiLK are telemetry/analytics integrations.

## O — Observability
Metrics, audit events and high-volume flow analytics are separated from transactional control state.

## P — PKI / Policy
Privileged mutations require authorization and, where policy requires, explicit approval.

## Q — Queueing
Bounded queues, connection reuse, worker pools and asynchronous telemetry prevent avoidable tail-latency amplification.

## R — Router
Registry-driven service endpoint routing with locality, health and measured path quality.

## S — Services
All FTN services integrate through stable contracts and service discovery rather than hard-coded router lists.

## T — Telemetry
ClickHouse is used for analytical workloads; operational state remains in the appropriate transactional store.

## U — Upgradeability
Provider/node additions are additive. Existing FamilyTimeNet DNS and service configuration must not be rebuilt merely to add a node.

## V — Verification
Every privileged change has an observed-state/verification step and an audit trail.

## W — Web / WebSocket
Web and WebSocket client/service boundaries remain provider-neutral and authenticated.

## X — eXternal Providers
External DNS/WAN/gateway providers are adapters. Provider-specific APIs must not leak into core contracts.

## Y — Yield / Resilience
Failover, drain-before-removal, stale-path eviction and bounded retries are required for resilient operation.

## Z — Zero Trust
Default-deny service access, least privilege, authenticated identities, explicit capabilities and auditable mutations.

## Release gate

The repository is structurally complete when all versioned contracts are internally consistent and CI passes. Real-world deployment remains dependent on authorized infrastructure credentials, real routers, BGP peers, DNS provider accounts and measured WAN paths; those cannot honestly be simulated or marked production-verified from source code alone.
