# FTN-AI Release Readiness

FTN-AI is the orchestration and AI/control layer for the FTN ecosystem. Existing implementations in FTN_service- and redy remain authoritative where they already provide a capability; FTN-AI consumes them through versioned contracts instead of duplicating control planes.

## Closed architecture domains

- Agent runtime and approval boundary
- Kernel wrapper and structured tool execution
- Backend/control-plane contracts
- Global and local full mesh
- Aether-Core path orchestration
- Global end-to-end latency scoring
- GoBGP routing control
- GeoIP2 locality intelligence
- Ansible desired-state orchestration
- Anycast lifecycle
- PKI/mTLS identity
- Gateway/edge proxy and authenticated SOCKS5
- WireGuard, AmneziaWG, OpenVPN 3, GRE, SSLH, Shadowsocks, Hysteria2 adapters
- BATMAN-adv, Yggdrasil and CJDNS mesh adapters
- NetSA/SiLK/YAF/super_mediator telemetry
- ClickHouse analytics boundary
- familytimenet.com registry-driven DNS
- additive DNS/provider/node onboarding
- Jupyter/ipykernel integration boundary
- Linux/openSUSE and FreeBSD host targets

## Production invariants

1. Adding a provider, node or gateway is additive; it must not require rebuilding familytimenet.com DNS.
2. No FTN-owned fiber or wireless backbone is assumed; Internet/WAN provider paths are transport resources.
3. DNS is first-class, but global traffic optimization is end-to-end rather than DNS-only.
4. Routing control state, FIB state, orchestration state and observed telemetry remain separate.
5. Privileged infrastructure changes require identity, capability/policy checks, approval where configured, verification and audit.
6. Private keys and credentials never belong in repository configuration.
7. SOCKS5/proxy gateways must not be unrestricted open relays.
8. PPTP is legacy-only and disabled by default.
9. Latency optimization must not bypass authentication or security controls.
10. RAM-resident/low-copy paths are used only when benchmark evidence supports them.

## Definition of done

A component is production-ready only after configuration validation, Go formatting/tests, security/policy tests, health/metrics contracts, failure behavior, audit behavior, idempotent registration, drain-before-removal for traffic nodes, and deployment documentation are present.

Environment-dependent BGP, router, DNS-provider, VPN and WAN integration tests require an authorized lab or production environment and are not falsely represented as unit-tested by repository CI.
