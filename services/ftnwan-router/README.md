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

## Execution boundary

The router control plane never accepts arbitrary shell commands, credentials, or AI-generated infrastructure mutations. Changes are represented as validated plans and passed through identity, policy, approval, preflight, transaction, verification and rollback controls.

## Dataplane boundary

VPP and DPDK are acceleration backends, not separate control planes. FTNWAN consumes normalized interface/route state through the Go contracts in `contracts/router.go`.

## MikroTik boundary

MikroTik is integrated through a dedicated adapter. NPK artifacts are versioned/signed build outputs; production credentials and private signing material remain outside the repository.

## Completion gate

A router release is production-ready only when the selected dataplane, routing, access, filtering, tunnel, management and provisioning adapters pass integration tests on the target hardware/OS. This repository contract does not claim that unavailable hardware-dependent adapters have been executed in this chat.
