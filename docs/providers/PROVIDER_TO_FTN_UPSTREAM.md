# Provider -> FTN Upstream

This is the production boundary for an upstream/transit provider handing connectivity to FTN.

## Traffic path

```text
Provider / Transit
      |
 Ethernet / VLAN handoff
      |
 FTN edge router
      |
 eBGP dual-stack
      |
 GoBGP / FTN routing control-plane
      |
 RPKI + bogon + prefix-length policy
      |
 FTN core / POP mesh
      |
 Customer / service networks
```

## Required provider data

Before a provider is admitted, FTN records:

- legal/canonical provider identity
- provider ASN
- FTN local ASN
- IPv4/IPv6 BGP neighbor addresses
- interface/VLAN handoff
- provider-advertised prefixes
- prefixes FTN will advertise
- default-route policy
- maximum accepted prefixes
- BFD capability and maintenance behavior
- NOC and abuse contacts
- committed capacity and interface speed
- maintenance window

Provider identity is not inferred from an IP address or ASN alone.

## Import policy

The default is reject. FTN only accepts explicitly permitted routes and applies:

1. RPKI-invalid rejection
2. bogon/private-prefix rejection
3. prefix-length limits
4. maximum-prefix protection
5. explicit default-route policy
6. community normalization/preservation
7. health-aware withdrawal

Customer routes must never be learned from an upstream provider through an implicit catch-all policy.

## Export policy

The default is reject. FTN exports only explicitly authorized FTN-originated prefixes. A provider peer cannot become an unintended transit path for customer, private, or management networks.

## Failover

Multiple providers are represented as independent eBGP peers. Path selection uses health, policy, local preference, and route validity. A degraded or failed provider may be withdrawn from selection, but any route mutation remains behind the FTN approval boundary.

## Operational controls

- BFD when supported
- graceful restart
- maximum-prefix protection
- RPKI validation when available
- route-change audit events
- pre-change snapshot
- post-change verification
- safe rollback
- peer state, prefix, traffic, latency and loss telemetry

## FTN implementation

The canonical GoBGP contract already declares eBGP/iBGP, multipath, graceful restart, BMP and FlowSpec capabilities while requiring approval for route changes. The provider contract in `configs/routing/provider-upstream-contract.yaml` adds the upstream-specific admission, import/export, health and safety boundary.

A provider remains `pending-admission` until its identity, policy and operational evidence have been reviewed. No provider record in this document grants permission to mutate live routing.
