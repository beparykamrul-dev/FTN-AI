# FTN Connectivity Stack v1

## Scope

FTN exposes one provider-neutral connectivity control plane. Protocols are adapters behind capability contracts; they are not independently managed applications.

## Transport adapters

| Adapter | Role | Production posture |
|---|---|---|
| WireGuard | Primary encrypted tunnel | Preferred |
| AmneziaWG | WireGuard-compatible specialized transport | Optional |
| OpenVPN 3 | VPN compatibility | Optional |
| GRE | Routed tunnel | Optional; must be policy restricted |
| SSL-VPN / SSLH | Compatibility/multiplexing edge | Optional |
| Shadowsocks | Proxy compatibility | Optional |
| Hysteria 2 | QUIC-based transport | Optional |
| PPTP | Legacy compatibility | Disabled by default; isolated |
| Socket/WebSocket | Application transport/control channel | API-dependent |

## Mesh adapters

- BATMAN-adv for supported Layer-2 Linux deployments.
- Yggdrasil and CJDNS as optional overlay-network adapters.
- Full-mesh topology is represented by the FTN mesh model rather than hard-coded protocol assumptions.

## DNS

`familytimenet.com` is a first-class FTN DNS namespace. DNS provider/server adapters remain behind the existing DNS abstraction. DNS operations require capability checks, audit records, and approval for privileged changes.

## PKI

All TLS/VPN identities use the FTN PKI abstraction. Private keys remain server-side. Certificate issuance, rotation and revocation are lifecycle operations, not arbitrary shell commands.

## Flow telemetry

NetSA components such as SiLK, YAF and super_mediator are treated as telemetry/analytics adapters. Normalized flow events are stored in ClickHouse for high-volume analytical workloads; transactional control state remains in the appropriate operational database.

## Development

Development environments must be reproducible. Jupyter kernels are registered from an explicit project kernel directory rather than a hard-coded `/path/to/...` value. `devenv` configuration may be used where available without making it a runtime dependency.

## OS targets

Linux/openSUSE and FreeBSD are supported host targets where individual adapters provide native support. Unsupported protocol combinations must fail capability negotiation instead of silently degrading.

## Security boundary

Every adapter declares capabilities, validates structured inputs, emits an audit event, and passes privileged operations through the FTN approval policy. No adapter accepts arbitrary client-supplied shell commands. Secrets are loaded from the server-side secret boundary and are never committed to the repository.

## Canonical lifecycle

`request -> capability check -> policy -> approval (if required) -> adapter -> observed state -> verification -> audit`
