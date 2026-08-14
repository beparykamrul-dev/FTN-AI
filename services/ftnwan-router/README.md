# FTNWAN Core Router

Production router subsystem boundary for FTN.OS / FTNWAN.

## Scope

- provider-neutral packet dataplane contract
- Linux kernel, VPP and DPDK integration points
- IPv4/IPv6 routing
- BGP/iBGP/GTSM control-plane integration
- VLAN and PPPoE access
- NAT, QoS and conntrack integration
- TCP/UDP and HTTP filtering boundaries
- VPN/tunnel and mesh adapters
- SSH/Telnet/router-management adapters
- MikroTik adapter and versioned NPK build pipeline
- device capability discovery and zero-touch provisioning
- policy-gated configuration transactions
- health verification and rollback
- FTN Metrics, audit and reconciliation integration
- shared AI assistant contract for local router, Android, PC, VPN and mesh access

## Security boundary

The router control plane does not accept arbitrary shell commands, credentials, or unrestricted AI-generated infrastructure mutations. Changes are represented as validated plans and passed through identity, policy, authorization, preflight, transaction, verification and rollback controls.

## Dataplane boundary

VPP and DPDK are acceleration backends, not separate control planes. FTNWAN consumes normalized interface/route state through the Go contracts in `contracts/router.go`.

## MikroTik boundary

MikroTik is integrated through a dedicated adapter. NPK artifacts are versioned/signed build outputs; production credentials and private signing material remain outside the repository.

## Access AI

Local router, Android, PC, VPN and mesh access surfaces share a bounded assistant contract. The assistant can explain, diagnose and recommend; privileged mutations remain behind authorization and transaction controls.

## CCTV

CCTV/cloud-player functionality is intentionally deferred until the core router and access/tunnel stack is complete.

## Completion gate

A router release is production-ready only when the selected dataplane, routing, access, filtering, tunnel, management and provisioning adapters pass integration tests on the target hardware/OS. This repository contract does not claim that unavailable hardware-dependent adapters have been executed in this chat.
