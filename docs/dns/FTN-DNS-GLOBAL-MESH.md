# FTN DNS — Global Mesh Architecture

FTN DNS is designed as a distributed DNS platform rather than a single recursive resolver.

## Layers

```text
Clients / PPPoE / VLAN / POP
          |
      Anycast DNS VIP
          |
       dnsdist
          |
   +------+------+
   |             |
 DNS Guard     Cache
   |             |
   +------+------+
          |
       Unbound
          |
  Adaptive Upstream Pool
          |
 Cloudflare / Google / Quad9 / AdGuard / OpenDNS / future providers
```

## Global mesh

Each FTN DNS location is a node in the global mesh. Anycast announces the same service address from healthy locations. Health state controls route advertisement/withdrawal so an unhealthy DNS site is removed from service.

The design separates:

- **Authoritative DNS:** PowerDNS primary with Knot secondary nodes.
- **Recursive DNS:** Unbound.
- **Traffic steering:** dnsdist.
- **Policy enforcement:** FTN DNS Guard.
- **Anycast:** BGP-based service advertisement from FTN DNS locations.
- **Geo/latency steering:** nearest healthy node and lowest-latency upstream selection.
- **Upstream resilience:** multiple independently operated DNS providers.

## Provider extensibility

Providers are registry entries, not hard-coded resolver logic. A provider can advertise one or more transports (UDP/TCP/DoT/DoH/DoQ), health probes, latency data, and priority. The resolver pool can therefore add future providers without redesigning the FTN DNS control plane.

## Performance path

1. Answer from local cache when valid.
2. Use stale cache during controlled upstream failure where policy permits.
3. Select a healthy low-latency upstream.
4. Use automatic failover when the selected upstream fails.
5. Keep per-node latency/error metrics for routing decisions.

## Security

DNSSEC validation is enabled for recursive resolution. Management APIs use mTLS. Raw DNS query logging is disabled by default; policy telemetry stores privacy-preserving domain hashes rather than raw domains.

## Production requirement

Anycast is only considered live after BGP advertisement, route health, DNS correctness, failover, cache behavior, and multi-location recovery are verified in the actual FTN environment. Configuration committed to Git is not evidence that the service is already globally live.
