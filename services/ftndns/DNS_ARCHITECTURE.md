# FTN DNS Architecture

## DNS identities

- Primary FTN DNS namespace: `familytimenet.com`
- Secondary/dynamic FTNDNS namespace: `ftndns.net`
- Dynamic service identity: `FTNDNS`

## Separation

`familytimenet.com` remains the primary FTN service namespace. `ftndns.net` is the dedicated dynamic DNS namespace for FTN-managed nodes and endpoints.

Dynamic records must be authenticated, validated and policy-controlled before authoritative state is changed.

## Resilience

FTNDNS should support multiple authoritative nodes, health-aware endpoint selection, controlled TTLs and reconciliation from the FTN main control plane. No single local network or POP should be a mandatory dependency for global DNS state.
