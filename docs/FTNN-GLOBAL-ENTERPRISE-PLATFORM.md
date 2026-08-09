# FamilyTimeNet Global Enterprise Platform

## Scope

A self-managed, enterprise control plane for authorized DNS, network, server, cloud, application, observability, data, and AI operations.

## Control Plane

- ASP.NET Core control/API plane
- Go backend and agents
- PostgreSQL system-of-record
- Redis/event/cache layer
- REST API
- GraphQL API
- WebSocket live operations
- Webhooks
- RBAC, multi-tenant and multi-workspace isolation
- SSO / LDAP / OAuth2 / OIDC
- Audit and approval engine

## DNS and Network Fabric

- PowerDNS
- Technitium DNS
- CoreDNS
- Unbound
- dnsdist
- GoDNS
- Anycast DNS
- Global DNS Mesh
- Tencent DNSPod
- Cloudflare DNS
- Akamai DNS
- DNS synchronization and health-aware failover

## Tencent Cloud Adapters

- DNSPod / DNS API
- COS
- CVM
- CLB
- VPC
- CDN
- Cloud Monitor
- SSL Certificate Manager
- IAM
- Anycast where available
- DNS synchronization
- Backup

Adapters MUST use provider APIs and least-privilege credentials. Credentials are never stored in source control.

## Infrastructure Graph

The canonical model represents:

- Sites, POPs and data centers
- Servers and VMs
- Routers, switches, OLT/ONU and network interfaces
- IPAM, prefixes, addresses and VLANs
- BGP sessions and circuits
- DNS zones and records
- Applications, services and containers
- Kubernetes resources
- Storage and cloud resources
- Dependencies and relationships

The graph separates **desired state** from **observed state**.

```text
Desired State
     |
Infrastructure Graph
     |
Live Discovery
     |
Reconciliation
     |
+----+-------------------+
| drift | impact | policy |
+----+-------------------+
          |
   Approved remediation
```

## Operational Intelligence

The web console provides authorized operators with:

- Global topology
- GIS/map view
- DNS and Anycast health
- BGP and network state
- Server and application health
- Cloud resources
- IPAM/DCIM/circuit inventory
- Dependency graph
- Drift and compliance findings
- Logs, metrics and events
- Incident and audit views

## Remote Operations

Remote operations are performed through authenticated agents and a controlled command router.

```text
Web Console
    |
WebSocket / REST / GraphQL
    |
Auth -> RBAC -> Policy -> Approval -> Audit
    |
mTLS Agent / Provider Adapter
    |
Server / Network / Cloud / Application
```

Operations are capability-based and allowlisted. Arbitrary unaudited shell execution is not part of the control plane.

High-impact operations such as routing, BGP, firewall, VLAN, DNS writes, credential changes, destructive deployment, and reboot require explicit authorization according to policy.

## AI Platform

- AI Dashboard Generator
- AI SQL Generator
- AI SQL Optimizer
- AI Visualization Recommendations
- AI KPI Detection
- AI Forecasting
- AI Insights
- AI Data Cleaning
- AI Anomaly Detection
- AI Auto Dashboard Builder
- AI Chat with Database
- AI Code Generator
- AI Report Generator

## Application Generation

Supported targets are planned as adapters/templates rather than hard-coded generators:

- PHP
- HTML
- React
- Vue
- Angular
- Flutter
- WordPress
- Laravel
- Joomla

## Data and Reporting

- PDF
- Excel
- CSV
- JSON
- XML
- REST
- GraphQL
- Webhooks
- 25+ connector architecture
- Scheduled and event-driven refresh

## Deployment

- Docker
- Kubernetes
- On-premise
- Self-managed infrastructure
- Air-gapped deployment profile
- Backup and restore
- Versioned configuration

## Security Model

- Encryption in transit with TLS/mTLS
- Encryption at rest where supported by the deployment
- Least privilege
- Secret manager integration
- RBAC and tenant isolation
- Immutable audit events
- Approval workflows
- Network segmentation
- No credentials in source control

## Commercial Entitlements

License, commercial-use, OEM, white-label, redistribution, support, update, and lifetime terms are deployment entitlements and contractual configuration. They must not be represented as automatically granted by the source code.
