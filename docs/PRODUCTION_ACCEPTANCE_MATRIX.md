# FTN-AI Production Acceptance Matrix

This matrix separates repository implementation from evidence obtained on authorized FTN infrastructure.

## Control plane

- [x] Health/readiness endpoints
- [x] Authentication boundary
- [x] Approval-gated privileged execution
- [x] Durable job lifecycle and idempotency
- [x] Event journal and audit trail
- [x] Recovery validation in CI
- [ ] Multi-node API failover on target infrastructure

## Core routing

- [x] GoBGP controller contract
- [x] Route-change approval requirement
- [x] Pre-change snapshot and post-change verification policy
- [x] RPKI/BFD/ECMP/GR policy declarations
- [ ] Real eBGP/iBGP peer acceptance
- [ ] BFD convergence measurement
- [ ] Core-A/Core-B failover measurement
- [ ] Route withdrawal/re-advertisement acceptance

## Network adapters

- [x] Typed adapter boundary
- [x] Ownership and safety execution gate
- [x] Bounded execution timeout
- [x] Destructive/external-action rejection
- [ ] Authorized MikroTik RouterOS read/telemetry acceptance
- [ ] Authorized OLT/ONU telemetry acceptance
- [ ] SNMP/NetFlow/IPFIX live collector acceptance

## DNS

- [x] DNS Guard policy contract
- [x] DNS Guard decision compiler
- [x] Resolver mesh contract
- [x] Runtime listener validation
- [ ] Authoritative delegation acceptance
- [ ] DNSSEC validation/signing acceptance
- [ ] DoT/DoH/DoQ listener acceptance
- [ ] Anycast/BGP withdrawal and recovery acceptance
- [ ] Multi-site resolver failover acceptance

## Security

- [x] FTN-owned asset scope
- [x] Manual-by-default privileged changes
- [x] Bounded defensive containment policy
- [x] Audit requirement
- [ ] Wazuh live API ingestion acceptance
- [ ] Credential rotation acceptance using secret store
- [ ] Firewall/rate-limit verification on target edge

## Data and recovery

- [x] Migration registry and recovery gate
- [x] Backup/restore CI evidence
- [ ] Target-environment backup restore drill
- [ ] Database failover acceptance
- [ ] Retention and storage-capacity validation

## Release rule

A component is **implemented** when source/configuration and automated validation exist. It is **production-verified** only after authorized FTN infrastructure passes the corresponding acceptance test with recorded evidence. No simulated result may be substituted for live evidence.
