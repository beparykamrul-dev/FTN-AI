# FTN Global Traffic Optimization v2

## Goal

Minimize stable end-to-end latency for FTN traffic, not merely DNS lookup latency. Decisions must consider RTT, jitter, loss, congestion, queueing, path length, endpoint health, and capacity.

## Open-source integration policy

FTN may integrate mature open-source implementations through explicit adapter boundaries rather than copying provider-specific proprietary behavior. Candidate building blocks include routing, telemetry, transport, DNS, and network-automation projects whose licenses and security posture have been reviewed. Each dependency is pinned, audited, and isolated behind a stable FTN interface.

Open-source code is an implementation source, not permission to access third-party networks or bypass provider controls. FTN only programs infrastructure it is authorized to operate.

## Data-plane latency model

```text
Client -> Access POP -> Aether-Core -> Backbone -> Egress POP -> Destination
                 |             |             |
              health        path score    capacity
                 +-------------+-------------+
                               |
                         route decision
```

A path score can combine:

- RTT
- jitter
- packet loss
- queue depth / utilization
- path stability
- endpoint health
- regional policy

The score is advisory to the routing layer and must have hysteresis so small measurement changes do not cause route flapping.

## RAM-to-RAM optimization

RAM-to-RAM is valuable inside a host or between processes because it avoids disk I/O and can reduce serialization/copy overhead. It does **not** make a physical network hop zero-latency. Network traffic still crosses the kernel/NIC/virtual-switch/physical link unless the communication is actually intra-host.

Use RAM-first paths for:

- kernel/backend queues
- local IPC/shared-memory buffers where justified
- hot configuration and route state
- connection/session caches
- telemetry batching
- worker pools

For same-host communication, prefer a measured low-copy IPC mechanism where it materially improves the workload. Do not add shared-memory complexity to ordinary control-plane calls unless profiling demonstrates a benefit.

For network traffic, optimize the actual data path: persistent connections, batching, appropriate transport, kernel/NIC tuning, queue management, local POP selection, and congestion-aware routing.

## Host and NIC optimization

The implementation should support measured tuning rather than unconditional system changes:

- CPU affinity for latency-sensitive workers where justified
- IRQ/RPS/RFS tuning where supported and measured
- NIC offloads only when benchmarked for the workload
- socket buffer sizing from measured bandwidth-delay product
- TCP pacing / modern congestion control where supported
- QUIC for workloads that benefit from it
- avoid unnecessary proxy hops
- reuse connections and avoid repeated handshakes
- bounded queues to prevent bufferbloat

Any host tuning is capability-detected and reversible.

## Global traffic

The controller should maintain a continuously refreshed topology view from active probes and passive telemetry. It can select an alternate POP or backbone path when a route becomes unhealthy or congested. It should not chase every transient RTT fluctuation.

DNS steering is only one input. Application and tunnel traffic must use the same global health/path model where the protocol supports it.

## Verification

Track before/after p50, p95, p99 RTT, jitter, loss, retransmissions, throughput, CPU, and queue depth. A change is considered beneficial only if tail latency and reliability improve without unacceptable capacity or security regressions.

## Security

Latency optimization never disables mTLS, authorization, PKI, audit, encryption, or approval boundaries. Legacy PPTP remains disabled by default.
