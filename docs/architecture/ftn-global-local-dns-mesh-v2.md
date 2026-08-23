# FTN Global + Local Full-Mesh DNS Architecture v2

## Design goal

FamilyTimeNet DNS is one logical service (`familytimenet.com`) implemented by independently deployable DNS nodes and provider adapters. Adding a provider or node is a registry/configuration operation, not a rewrite of the existing DNS service.

## Full mesh

Both local POPs and global POPs participate in the same provider-neutral mesh model. The mesh stores node/link health and measured RTT, loss and jitter. Route selection can prefer locality while retaining global failover.

## Routing stack

- GoBGP is the BGP control-plane integration rather than a generic BGP-only abstraction.
- GeoIP2 provides geographic classification for policy and locality decisions.
- Aether-Core supplies the provider-neutral path/orchestration layer.
- DNS remains an independent service plane and is not coupled to BGP implementation details.

## DNS contract

Every DNS provider/node implements the same FTN DNS v1 contract. The registry stores provider identity, scope, region, endpoint, zone and capabilities. Provider-specific SDKs remain behind adapters.

The canonical zone is `familytimenet.com`. Existing nodes continue serving the zone when a new node is added. A new provider is registered, health-checked, synchronized and then eligible for traffic according to policy.

## Additive onboarding lifecycle

`register node -> validate contract -> health check -> zone consistency check -> warm/cache -> advertise eligibility -> serve traffic`

Removal uses a drain lifecycle rather than deleting the canonical DNS model first.

## Local versus global

Local mesh:
- low-latency POP-to-POP connectivity
- local DNS resolution/cache
- local failure isolation

Global mesh:
- inter-region connectivity
- global failover
- latency-aware path selection
- Anycast/BGP integration where appropriate

The two scopes share the node and telemetry model but retain separate policy controls.

## Provider extensibility

The existing PowerDNS/Technitium/CoreDNS/Unbound/dnsdist/GoDNS/Anycast/DNSPod/Cloudflare/Akamai integrations can coexist with FTN-owned nodes. New providers should only implement the adapter contract and register a node/config entry.

No provider is allowed to mutate another provider's credentials or implementation state.

## Verification

A new node is eligible only after DNS health, zone consistency, endpoint health and mesh path checks pass. Route changes are measured using p50/p95/p99 latency, loss and jitter before and after activation.
