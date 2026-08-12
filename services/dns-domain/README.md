# FTN Domain + DNS Panel

Production service boundary for the FTN control panel's domain and DNS management.

## Panel capabilities

- Domain inventory and status
- Zone/record inspection
- A / AAAA / CNAME / TXT / MX / NS / SRV / CAA records
- TTL and priority management
- DNSSEC status
- Nameserver visibility
- Resolver health integration
- Propagation checks
- Change audit

## Safety

All mutations must pass FTN identity, RBAC, policy, approval, backend validation and audit controls. This package contains contracts only; provider-specific mutation adapters must implement the contract separately.
