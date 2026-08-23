# FTN Global Latency + Aether-Core v1

## Objective

Latency optimization applies to the complete FTN network path, not only DNS. The control plane continuously evaluates POP-to-POP and edge-to-service paths using RTT, jitter, packet loss, utilization and health.

## Aether-Core role

Aether-Core is represented as a provider-neutral latency-aware fabric layer. It does not replace WireGuard, routing protocols, DNS, or mesh implementations. Instead it supplies common path-health and selection policy to those adapters.

```text
Client
  |
Nearest healthy ingress
  |
Aether-Core path policy
  |
+------------------------------+
| POP / Backbone / Mesh        |
| WG / AMNEZIAWG / Hysteria2   |
| BATMAN-adv / Yggdrasil/CJDNS |
+------------------------------+
  |
Nearest healthy egress/service
  |
Application / DNS / API / Flow
```

## What is optimized

- DNS lookup latency
- TCP/QUIC connection establishment
- VPN tunnel path latency
- application/API latency
- POP-to-POP backbone latency
- inter-region traffic
- mesh peer selection
- congestion and queueing delay
- packet loss and retransmission effects
- jitter and tail latency

DNS is therefore one consumer of the global path system, not the system itself.

## Path scoring

The first implementation uses a deterministic score:

`score = RTT + 0.35*jitter + 1000*loss + 200*utilization`

Lower is preferred. These coefficients are policy defaults, not universal network constants. They must be validated against FTN measurements before production tuning.

## Control-loop

`probe -> normalize -> score -> select -> route -> observe -> verify -> retain/rollback`

Hysteresis is required so small measurement changes do not cause route flapping. Stale or unhealthy paths are excluded. Regional and global failover are separate decisions.

## Data sources

Active probes provide RTT/jitter/loss. Passive telemetry provides utilization and flow evidence. NetSA components (SiLK/YAF/super_mediator) feed normalized analytics into ClickHouse. Operational state remains in the control-plane database.

## Transport integration

The policy is transport-neutral and can select paths carrying WireGuard, AmneziaWG, OpenVPN 3, Hysteria2, WebSocket or other approved FTN transports. It does not weaken authentication or encryption to reduce latency.

## Mesh integration

BATMAN-adv, Yggdrasil and CJDNS peers expose health/latency observations to the same selection layer where supported. Full-mesh does not mean every packet must use the shortest geographic route; the measured end-to-end path wins subject to policy, capacity and security constraints.

## Performance controls

- persistent connections and keep-alive
- warm workers and kernel tool registry
- connection pooling
- parallel probing
- asynchronous telemetry
- bounded queues
- local POP preference with measured-path override
- fast unhealthy-path eviction
- controlled hysteresis to prevent route oscillation

## Safety

Aether-Core is a control-plane abstraction, not an unrestricted packet manipulation interface. Privileged changes remain behind FTN authentication, capability checks and approval policy. PPTP remains disabled as a latency fallback.
