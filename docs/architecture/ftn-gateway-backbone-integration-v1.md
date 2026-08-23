# FTN Gateway Backbone Integration v1

## Purpose

FTN-AI is the software control/data-plane for an ISP architecture that does not assume an FTN-owned fiber or wireless backbone. Provider and Internet WAN paths are therefore first-class transport resources.

## Existing repository alignment

The existing `FTN_service-` repository is an implementation source, not a second control plane. Its platform configuration already declares an Internet-WAN-only backbone, centralized control/telemetry, DNS mesh/Anycast, PowerDNS/Technitium/CoreDNS/Unbound/dnsdist/GoDNS and external DNS providers. FTN-AI must consume and normalize these capabilities rather than recreate or fork them.

The `redy` repository is the broader enterprise foundation covering kernel/core lifecycle, unified control plane, socket/event fabric, network automation, DNS/global mesh, monitoring, PKI/RBAC/audit and deployment. FTN-AI integrates these contracts where they already exist.

## Gateway fabric

All approved gateway types implement the same FTN Gateway contract:

- routing table / policy routing
- GoBGP control-plane integration
- Ansible orchestration
- Anycast advertisement/withdrawal
- PKI identity and mTLS
- edge/reverse proxy
- HTTP/HTTPS proxy adapters
- HEDSHOW/other project-specific gateway adapters when an implementation contract exists
- SOCKS5 gateway
- VPN/tunnel gateways (WireGuard, AmneziaWG, OpenVPN 3, Hysteria2 and other supported transports)

Gateway selection is based on capability, health, route availability, latency, loss, utilization, policy and scope. A gateway is not automatically trusted because it exists in the registry.

## Routing

The routing model separates:

1. RIB/control state (GoBGP and provider routing)
2. FIB/kernel state (host routing tables/policy routing)
3. gateway service state (proxy/VPN/Anycast edge)
4. observed path telemetry

Ansible is an orchestration mechanism, not the routing source of truth. Desired state comes from the FTN control plane and is rendered into idempotent playbooks/roles.

## Anycast

Anycast nodes advertise only after:

`identity -> capability -> health -> route validation -> policy -> advertisement`

Withdrawal is triggered by failed health/SLO checks, route loss or policy state. Hysteresis prevents flapping.

## PKI

Gateway identities use the FTN PKI contract. Private keys remain in the server-side secret boundary. Certificates are rotated and revoked through lifecycle operations.

## Edge proxy

The edge proxy layer is provider-neutral. It can terminate TLS, route by hostname/path, forward to an approved upstream, and emit telemetry. Proxy configuration is generated from the same gateway contract; provider-specific configuration is an adapter detail.

## SOCKS5

SOCKS5 is an explicitly registered gateway capability. It is subject to authentication, authorization, egress policy, rate limits, audit and abuse controls. It is not an unrestricted open relay.

## Global + local full mesh

The same gateway contract applies to local and global nodes. Local mesh peers are preferred when they meet the path policy; global peers provide alternate routes. Aether-Core supplies path scoring/orchestration above the transport implementations.

## Latency

Path selection considers RTT, jitter, packet loss, utilization, queueing and stability. Hot routing state and connection pools may remain memory-resident. RAM-to-RAM/low-copy paths are used only where benchmarks demonstrate benefit; physical network propagation and switching latency remain irreducible.

## Additive integration invariant

Adding a provider, gateway, DNS node or POP must be additive:

`registry entry -> contract validation -> capability registration -> health -> synchronization -> eligible`

Existing FTN services are not rebuilt merely because a new provider/node is introduced.
