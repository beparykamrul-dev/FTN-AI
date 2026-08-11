# FTN Service Platform

FTN is organized as one governed platform with pluggable services. The Control Plane is the stable boundary; services consume canonical infrastructure, identity, policy, event, and audit APIs.

## Service planes

- **Control Plane:** administration, RBAC, tenants, approvals, audit, configuration, deployments, WebSocket events.
- **Infrastructure Graph:** devices, links, services, IPs, circuits, fiber, OLT/ONU, customers and relationships.
- **Discovery & Assurance:** discovery, normalization, desired-vs-observed reconciliation, drift and compliance.
- **Network Plane:** Core/POP/Client/Backup nodes, routing, BGP, Anycast, DNS, firewall, WireGuard, QoS, eBPF/XDP.
- **Mesh Plane:** node/link health, topology, path selection, failover proposals and convergence orchestration.
- **Observability Plane:** metrics, logs, traces, SNMP, flow telemetry, syslog, eBPF and network inspection.
- **Service Plane:** billing, customer services, IPTV, CCTV, cache, DNS, portals and partner integrations.
- **Builder Plane:** source import, package/dependency detection, build, test, image/package generation and governed deployment.
- **Client Plane:** Android, PC and Smart TV applications with FTN identity and service policy.
- **AI Plane:** grounded analysis, anomaly detection, impact analysis, recovery proposals and post-change verification.

## Node products

FTN Core OS, FTN POP OS, FTN Client OS and FTN Backup OS are deployment profiles of the same FTN platform, not independent control systems.

## Governance invariant

All production writes follow: identity -> authorization -> policy validation -> change set -> approval (when required) -> execution -> verification -> audit. Existing configuration is backed up before destructive replacement, and rollback metadata is retained.

## API/event boundary

REST/GraphQL APIs are used for request/response operations. WebSocket/event streams are used for live state, topology and job updates. External automation and AI integrations enter through governed APIs/MCP rather than bypassing the Control Plane.

## Customer and infrastructure relationship

Customer, subscription, PPPoE, IP, ONU, OLT, VLAN, fiber, router and service records are linked through the Infrastructure Graph with tenant/RBAC boundaries. Customer-facing clients receive only authorized records and actions.
