# FTN Kernel + Backend + Backbone v1

## Purpose

Unify the FTN network kernel, backend control plane, and network backbone under one provider-neutral lifecycle.

```text
Jupyter/ipykernel
      |
FTN Kernel Wrapper
      |
Capability + Approval Policy
      |
Backend API / Control Plane
      |
Job Queue + Event Journal
      |
Device Adapter Layer
      |
Network Backbone
  |       |       |
VPN     DNS     Mesh
  |       |       |
Telemetry / NetSA / ClickHouse
```

## Kernel boundary

The kernel is a controlled orchestration client. `do_execute` accepts structured FTN tool requests and returns structured results. It must not become an unrestricted Python or shell gateway.

## Backend boundary

The backend owns:

- authentication and mTLS identity
- enrollment and revocation
- RBAC/capability checks
- approval state
- idempotent jobs and leases
- audit/event journal
- inventory and desired state
- adapter dispatch
- verification and observed state

Existing control-plane authorization remains authoritative; the kernel cannot bypass it.

## Backbone boundary

The backbone is represented as provider-neutral resources:

- encrypted transports: WireGuard, AmneziaWG, OpenVPN 3
- routed/compatibility transports: GRE, SSLH, Shadowsocks, Hysteria2
- application transport: Socket/WebSocket
- legacy compatibility: PPTP disabled by default
- mesh: BATMAN-adv, Yggdrasil, CJDNS
- DNS: familytimenet.com and DNS provider adapters
- PKI: service/device identity and certificate lifecycle
- telemetry: NetSA/SiLK/YAF/super_mediator and ClickHouse

## Execution contract

`request -> authenticate -> authorize -> approve -> enqueue -> lease -> adapter -> observe -> verify -> journal`

Read-only discovery may be automated when the policy permits it. State-changing network operations require the configured approval policy. Every state-changing operation must be idempotent where possible and produce an audit event.

## Dependency isolation

Python tools such as Netmiko, Paramiko, NAPALM, PySiLK and ipykernel are optional worker/runtime integrations. The Go backend remains the system of record and does not import Python libraries directly. A Python worker communicates through the structured FTN tool contract.

## Backpressure and resilience

The backend must bound concurrent jobs, use leases for work ownership, retry only idempotent operations, preserve event ordering, and expose health/readiness information. Device failures must not block unrelated tenants or devices.

## Security

- mTLS for trusted agents
- server-side secrets only
- no arbitrary shell commands from API input
- no client-provided private keys
- explicit capability allowlists
- approval before privileged state changes
- immutable/auditable operation identifiers
- verification after mutation

## Data plane versus control plane

Control-plane state belongs in the operational database. High-volume flow/telemetry belongs in ClickHouse. Large artifacts belong in object storage. The kernel never becomes a database or source of truth.
