# FTN-AI Final Release

## Release status

FTN-AI source tree is frozen at the architecture/contract level. Future changes must be additive and pass CI; they must not bypass the registry, policy, security, audit, verification, or service-contract boundaries.

## Included

- Kernel / ipykernel structured tool boundary
- Backend control-plane contracts
- Registry-driven FTN service router
- FamilyTimeNet DNS registry for `familytimenet.com`
- Additive DNS/provider/node onboarding
- Local and global full-mesh model
- Aether-Core path-selection model
- GoBGP + GeoIP2 routing integration boundary
- Anycast lifecycle
- Ansible desired-state orchestration boundary
- PKI/mTLS and approval/audit boundaries
- Edge/proxy/SOCKS5 gateway contracts
- WireGuard, AmneziaWG, OpenVPN 3, GRE, SSLH, Shadowsocks, Hysteria2 and WebSocket integration boundaries
- BATMAN-adv, Yggdrasil and CJDNS mesh adapters
- NetSA/SiLK/YAF/super_mediator telemetry boundary
- ClickHouse analytical telemetry boundary
- End-to-end latency model using RTT, jitter, loss, utilization, queue depth and stability
- Router and DNS registry regression tests
- Repository and architecture-contract CI gates
- A-Z completion and deployment-acceptance documentation

## Non-fake production gate

Source completeness does not manufacture live infrastructure evidence. Production acceptance requires real authorized FTN infrastructure, including DNS service endpoints, routers, GoBGP peers, gateway nodes, PKI services and WAN paths. No simulated health, route, latency or provider result is treated as production verification.

## Change policy

1. Register new providers/nodes/services instead of rebuilding existing FamilyTimeNet DNS.
2. Keep provider-specific implementations behind adapters.
3. Keep routing control, forwarding state and service routing separate.
4. Preserve approval, least privilege, mTLS, audit and post-change verification.
5. Measure end-to-end latency; do not optimize DNS alone.
6. Prefer measured improvements over speculative kernel/network tuning.
7. Run repository validation and architecture-contract validation for every main/PR change.
