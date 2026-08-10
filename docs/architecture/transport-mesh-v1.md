# FTN Transport & Mesh v1

FTN uses a provider-neutral transport abstraction. The following technologies are integration adapters/reference implementations, not mandatory dependencies.

## Secure mesh / tunnels

- WireGuard: primary low-level encrypted tunnel adapter.
- NetBird: optional mesh management/control adapter.
- Hysteria 2: optional QUIC-based transport adapter where appropriate.
- Keepalive: health/liveness contract for every FTN agent and mesh link.
- BATMAN-adv: optional Layer-2 mesh adapter for supported Linux wireless/edge deployments.
- Babel: optional dynamic routing adapter for mesh networks.

## Network/security edge

- OPNsense: firewall/router integration adapter.
- FD.io VPP: high-performance dataplane adapter.
- OpenVPN Enterprise: compatibility VPN adapter.
- SSL-VPN: generic compatibility adapter; vendor-specific products remain separate adapters.
- Palo Alto GlobalProtect and Fortinet FortiSSL VPN: optional enterprise VPN integration adapters, subject to vendor APIs/licensing.
- CrowdSec: security signal/remediation adapter.

## Data movement / local transfer

- Resilio Sync: optional managed file synchronization adapter.
- Wormhole / Magic Wormhole: optional authenticated point-to-point transfer adapter.
- USB automount: FTN host/device policy module for approved removable media workflows.
- Thunderbolt: host hardware/interface inventory and policy integration; it is not treated as a network transport.

## Commercial VPN provider

NordVPN is treated as an optional outbound VPN/provider integration, not as the FTN mesh core. Provider credentials must remain in the FTN secrets manager and never in source code.

## Canonical FTN flow

Device/Server -> FTN Agent -> authenticated transport -> Mesh Control Plane -> policy/RBAC/approval -> adapter -> dataplane -> observed state -> verification.

The Control Plane must not execute arbitrary shell commands from an untrusted client. Every adapter declares capabilities, validates inputs, and records an auditable operation. Production changes require explicit authorization according to the FTN approval policy.
