# FTN Unified Architecture V2

## Scope

FTN-AI is the unified control/orchestration layer for a modular, local-first and globally interoperable FTN platform. Production implementations remain in their source repositories and are consumed through explicit contracts and adapters.

## Control Plane

- ASP.NET Core control-plane APIs
- Go backend/services
- Central PostgreSQL with PgBouncer
- Central audit/tracking
- Approval-first state changes
- Stable service/node/provider identity and capability registry

## Core Network Fabric

- Aether-Core path/orchestration abstraction
- Multi-WAN path policy
- WireGuard and AmneziaWG adapters
- Experimental transports only behind explicit capability flags
- GoBGP/RIB control integration
- Kernel FIB and policy routing
- eBPF/XDP enforcement and telemetry hooks
- Automatic health-aware failover and reconciliation

## DNS Fabric

- FTN DNS and global DNS mesh
- PowerDNS
- Technitium
- CoreDNS
- Unbound
- dnsdist
- GoDNS
- Hickory DNS
- Anycast DNS
- External DNS provider adapters
- DNS API and management
- DNSSEC and certificate-aware DNS operations

## Security

- Private PKI and mTLS
- ACME certificate lifecycle
- OCSP/CRL lifecycle
- TPM 2.0/HSM integration boundaries
- RBAC and approval queue
- Secrets kept outside source control
- WAF/DDoS/bot/threat-intelligence integrations only through authorized providers
- Security auditing and device-trust signals

## Observability and Flow

- Metrics and health polling
- NetFlow/IPFIX/flow adapters
- SiLK/NetSA/YAF-compatible analytics
- Nfdump/NfSen/pmacct adapters
- ClickHouse for high-volume analytical telemetry
- Service and infrastructure audit trail
- Geo/GIS and provider/network map

## Hosting, Cloud, CDN and Edge

FTN is modeled as a future Cloud/CDN/Edge provider while retaining interoperability with external providers.

- FTN Hosting
- FTN Cloud
- FTN CDN
- FTN Edge/POP
- Cache services
- Object/file storage
- VPS/container hosting
- Static/web/application hosting
- External cloud/CDN/hosting provider adapters

Local services are first-class. Global providers are optional capacity/interoperability layers rather than mandatory dependencies.

## Application and Platform Services

- FTN Android client
- Windows/Linux/macOS integration targets
- FTN Drive
- FTNBOX
- Web Builder
- Android Builder
- Module Builder
- AI Agent
- Corporate/agent/reseller platforms
- Billing/accounting/inventory
- Employee management
- AI call center/reporting
- E-commerce/social/application services
- IPTV/media/CCTV/local-service modules

## Mesh and Service Fabric

Local and global nodes use the same identity/capability contract:

`register -> validate -> authenticate -> health -> sync -> route eligibility`

Service lifecycle:

`request -> authorize -> desired state -> reconcile -> execute -> health verify -> audit`

Adding a provider, node, POP or service must be additive and must not rebuild existing FTN DNS, routing, gateway or service state.

## Data Boundaries

- Transactional/source-of-truth state: PostgreSQL-backed service contracts
- High-volume analytical telemetry: ClickHouse/flow analytics
- Routing control: GoBGP/control plane
- Forwarding: kernel FIB/policy routing
- Path orchestration: Aether-Core abstraction
- DNS membership: canonical FTN DNS registry
- Credentials/secrets: runtime secret boundary

## Production Acceptance

A technology or repository is not considered production-integrated merely because it is named here. It becomes eligible only after a concrete adapter/contract, tests, security review, CI validation and successful acceptance on authorized infrastructure.
