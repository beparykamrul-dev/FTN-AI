# FTN Upstream, Edge and Monitoring

## Scope

The upstream adapter is the control-plane boundary for FTN ISP upstream, approved provider CDN/interconnect, and IX connections. It does not intercept or copy third-party encrypted content.

## Traffic model

Customer -> OLT/BNG -> FTN POP -> FTN Core -> approved interconnect/CDN or transit -> provider

Provider -> approved interconnect/CDN or transit -> FTN Border -> FTN Core -> POP -> BNG/OLT -> customer

CDN/cache is used only where the provider authorizes the interconnect/cache arrangement. General Internet traffic continues through the normal BGP/transit/peering path.

## Monitoring

Collect per-upstream health, latency, packet loss, prefix count/limit, last check time and state. Alert on down state, elevated latency/loss and route-table prefix exhaustion. Keep payloads and credentials out of telemetry.

## Control boundary

Configuration changes require the FTN approval policy. Automatic failover may select an already-approved healthy endpoint; adding peers, prefixes, credentials, or changing routing policy remains an administrative action.
