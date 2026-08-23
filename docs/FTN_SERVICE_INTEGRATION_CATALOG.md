# FTN Service Integration Catalog

## Purpose

FTN-AI is the orchestration/control layer; it must not duplicate production services that already live in the FTN service repositories. Existing implementations are integrated through explicit contracts, adapters and registry entries.

## Repository roles

### FTN_service-
Production-oriented FTN service foundation. Its documented scope includes the unified global service cluster, central PostgreSQL/PgBouncer, telemetry/audit, geographic topology, Aether-Core transport abstraction, global DNS mesh, ASP.NET Core/Go control-plane contracts and the DNS provider ecosystem. It also defines HTTP health/readiness endpoints and approval-first state changes.

Integration rule: FTN-AI consumes these capabilities through service/provider contracts and must not create a competing source of truth.

### FTNDNS_AI
DNS-focused source repository. Its current main branch contains source-map/import documentation rather than a duplicate runtime implementation. FTN-AI treats it as a DNS research/source input and keeps `familytimenet.com` registry semantics canonical.

### redy
Large multi-purpose foundation repository. FTN-AI integrates only the FTN-relevant, reviewed components/contracts from it; unrelated generated/example workflows are not promoted into the FTN production runtime.

### FTN_ser_AI
Currently has no README/runtime contract on its main branch. It is therefore registered as an integration candidate, not claimed as production-ready until a concrete service contract and implementation are present.

## Canonical service families

- DNS: PowerDNS, Technitium, CoreDNS, Unbound, dnsdist, GoDNS, Anycast, external provider adapters
- Routing: GoBGP, kernel FIB/policy routing, GeoIP2
- Gateway: Anycast edge, reverse proxy, authenticated SOCKS5, VPN/tunnel adapters
- Transport: WireGuard, AmneziaWG, Aether-Core adapters and approved transport implementations
- Network: router/device adapters, topology and health discovery
- Telemetry: metrics/events, flow analytics, NetSA/SiLK/YAF/super_mediator, ClickHouse analytical storage
- Security: PKI/mTLS, RBAC, approval queue and audit trail
- Application: FTN web/mobile/service endpoints through stable API contracts

## One registry rule

Every service/node/provider is represented by a stable identity and capability contract:

`register -> validate -> authenticate -> health -> sync -> route eligibility`

Adding a new implementation must be additive. Existing `familytimenet.com` DNS, routing, gateways or service configuration must not be rebuilt merely because a new node/provider is added.

## Source-of-truth boundaries

- Transactional platform state: PostgreSQL/source-of-truth services
- Analytical/high-volume telemetry: ClickHouse/flow analytics
- Route control: GoBGP/control plane
- Forwarding: kernel FIB/policy routing
- Service endpoint selection: FTN Router
- Global path orchestration: Aether-Core
- DNS node/provider membership: FamilyTimeNet DNS registry
- Credentials/secrets: runtime secret boundary; never source-controlled

## Production rule

No source repository is considered production-integrated merely because it is listed here. Integration requires a concrete contract, implementation, tests, security review and successful CI/live acceptance for the authorized target infrastructure.
