# FTN Low-Latency Architecture v1

The latency layer optimizes path selection and execution overhead without weakening the FTN security/control boundary.

## Runtime path

```text
Client
  -> nearest healthy FTN edge/POP
  -> persistent secure transport
  -> warm kernel/tool registry
  -> pooled backend connection
  -> adapter worker
  -> asynchronous telemetry
```

## Main latency controls

1. **Latency-aware routing** — retain recent endpoint health/RTT/loss samples and select the fastest healthy endpoint.
2. **Connection reuse** — keep persistent secure sessions where the protocol supports them; avoid repeated handshakes.
3. **DNS caching/prewarm** — cache positive answers according to authoritative TTL and prewarm only approved FTN destinations.
4. **Warm workers** — keep the kernel registry and Python network workers warm instead of spawning a process per request.
5. **Bounded concurrency** — avoid CPU/socket contention that increases tail latency.
6. **Local POP preference** — prefer a healthy nearby FTN POP, then fail over to the fastest healthy peer.
7. **Async telemetry** — telemetry and ClickHouse writes are decoupled from interactive control requests.
8. **Idempotent retry only** — retries use bounded exponential backoff and never turn a non-idempotent mutation into an unsafe automatic retry.
9. **QUIC/HTTP2 when supported** — use transport features that reduce connection setup and head-of-line effects when compatible with the selected service.

## Security invariants

Latency optimization never bypasses mTLS, enrollment, capability checks, approval, audit, verification, or secret isolation. PPTP is not a latency fallback. No route is selected solely because it is fast if it is unhealthy or violates policy.

## Measurement

The router consumes endpoint observations rather than raw packet payloads. At minimum track RTT, loss, health, timestamp, and selected endpoint. Measure p50/p95/p99 latency separately for DNS, kernel request, backend dispatch, adapter execution, and end-to-end network paths.
