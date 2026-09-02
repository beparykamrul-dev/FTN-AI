# FTN DNS — Install vs Portal Adjust

## Installed runtime components

These are host/container services and must be installed and health-checked:

- dnsdist — DNS frontend, routing and failover
- Unbound — recursive resolver, cache and DNSSEC validation
- CoreDNS — cloud/mesh DNS node
- Technitium DNS — optional Enterprise management/recursive node
- Caddy — optional ACME edge component
- PowerDNS Enterprise — licensed authoritative service; deploy its vendor package/image and connect it to FTN control plane

## Portal-managed components

These are configuration/data integrations exposed through the FTN DNS Portal rather than blindly installed as local daemons:

- FTN Global DNS Mesh topology
- Anycast node inventory and health state
- DNS zones, records, TTL and DNSSEC intent
- PowerDNS/Knot/Technitium/CoreDNS node metadata
- dnsdist pools and routing policy
- Unbound policy and approved upstream pool
- Cloudflare DNS
- Tencent Cloud DNSPod
- Akamai DNS
- DuckDNS
- Porkbun
- GoDNS automation
- Hickory DNS integration
- Numa DNS integration
- Fastly edge integration
- Let's Encrypt DNS-01 provider configuration
- Caddy DNS provider selection

## Safety and operational boundary

Provider API tokens, private keys, certificates and BGP credentials are supplied through the deployment secret store. They are never committed to this repository.

Anycast/BGP announcement is a runtime network operation. The portal may produce and validate an intent/configuration, but live advertisement requires the FTN router/BGP adapter and successful health gates.

## Target request path

Client → Anycast VIP → dnsdist → DNS policy/Guard → Unbound/CoreDNS → cache/DNSSEC → approved upstream mesh.

Authoritative requests use the FTN authoritative cluster and its secondary/replication path.
