# FTN-AI Production Readiness Checklist

This checklist separates repository implementation from live infrastructure acceptance.

## Control plane
- [x] Authenticated API boundary
- [x] Approval-gated privileged changes
- [x] Durable jobs and event journal
- [x] Idempotency and lease recovery
- [x] Audit trail

## Network control
- [x] Typed adapter contract
- [x] Pre-change snapshot requirement
- [x] Validation and post-change verification
- [x] Rollback-when-safe policy
- [x] No arbitrary remote shell execution

## Core routing
- [x] GoBGP controller contract
- [x] BGP/BFD/ECMP/GR policy
- [x] RPKI enforcement when available
- [x] Core HA/failover contract
- [ ] Live peer acceptance on FTN-owned routers
- [ ] Live failover/recovery acceptance

## DNS
- [x] DNS Guard policy contract
- [x] Resolver/orchestrator contract
- [x] Listener validation
- [ ] Live DoT/DoH/DoQ listener acceptance
- [ ] Authoritative delegation acceptance
- [ ] Multi-node Anycast/BGP acceptance

## Observability
- [x] Event and audit contracts
- [x] Recovery validation gates
- [ ] Live SNMP/NetFlow/IPFIX feeds
- [ ] Live alert routing and notification acceptance

## Access/OLT
- [x] Adapter contract
- [ ] Live MikroTik acceptance
- [ ] Live OLT/ONU vendor acceptance
- [ ] PPPoE/VLAN/IPAM end-to-end acceptance

## Release gate

A production release requires all applicable live acceptance items above to pass on authorized FTN-owned infrastructure. Repository code, contracts, or CI success alone must not be represented as proof of live BGP, DNS, router, OLT/ONU, WAN failover, or availability performance.
