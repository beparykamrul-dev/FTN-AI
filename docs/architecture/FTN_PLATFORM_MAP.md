# FTN Platform Map

## Core principle
FTN is designed as a multi-server, interconnected platform. Servers may be geographically or operationally separate, but communicate through authenticated service APIs and explicit network relationships.

## Platform domains

- FTN Global Friendly DNS / FTN DNS
- FTN Hosting / Global Hosted Services
- FTN Backend
- FTN Android Client
- FTN Router / Firewall / Web UI / Android Control
- FTN CCTV / CCTV Cloud
- FTN Self-Core Router
- FTN Mesh
- FTN Data Consistency
- FTN WAN / SDN
- FTN Server Fabric
- FTN ACME Key Manager
- FTN WebSocket / API / Proxy / Metrics
- FTN Device, Connection, IP, Relationship and DCIM management
- FTN Drive
- FTN Web Builder
- FTN Android Builder
- FTN Corporate Client Platform
- FTN Remote / FTNBOX local mode
- FTN AI Agent
- FTN Module Builder
- FTN TV App / Client / TV Server
- FTN Windows Services
- FTN Cache
- FTN Agent / Reseller Platform
- FTN Apple Services
- FTN CDN / Edge
- FTN DNS Firewall / SDK
- FTN Accounts / Billing / Inventory
- FTN Employee Platform
- FTN AI Analysis
- FTN E-commerce
- FTN AI Call Center / Reports
- FTN Geo Map / Fiber GIS
- FTN Go Services
- FTN Social App

## Server relationship model

Each server is an independently deployable node with:

1. unique service identity;
2. explicit allowed peers;
3. health and capacity reporting;
4. authenticated API communication;
5. local-first operation where practical;
6. controlled synchronization with FTN Data Consistency;
7. centralized policy without requiring a single point of failure.

The control plane must never require every data-plane node to depend on one server for normal operation.

## Security boundary

Sensitive services are separated from the public proxy/API path. Certificate private keys, payment credentials, identity secrets and other high-value secrets must not be exposed to ordinary request handlers.

Future cybersecurity capabilities are extension points rather than assumptions in the core data plane.

## Local-to-global model

FTN services can operate locally first and extend to global nodes through FTN WAN/SDN, mesh, DNS, CDN/edge and authenticated service APIs. Local service availability should degrade gracefully when an external provider is unavailable.

## Product roadmap boundary

This document records the platform architecture and intended service boundaries. Commercial availability, regulatory registration, licensing and third-party service eligibility must be validated separately before offering a service publicly.
