# FTN Router → Service Fabric v1

## Goal

Provide one routing contract from an FTN router/gateway to every registered FTN service. The router is provider-neutral and does not require rebuilding services when a new node or provider is added.

## Service classes

The registry can represent DNS, control-plane APIs, kernel workers, monitoring, telemetry, PKI, proxy/edge, VPN/tunnel gateways, mesh services, billing/application services and future FTN services through the same endpoint contract.

## Resolution

`service -> healthy candidates -> locality preference -> RTT/loss/priority score -> gateway -> endpoint`

The hot result is cached briefly in memory. Cache expiry is bounded so topology and health changes are observed quickly. A failed lookup returns an explicit no-route error rather than silently selecting an unhealthy endpoint.

## Router/Gateway boundary

Routing control (GoBGP), host FIB/policy routing, Anycast advertisement, Aether-Core path scoring and service endpoint resolution remain separate concerns. They exchange versioned state rather than sharing mutable implementation internals.

## Failure handling

- unhealthy endpoints are excluded;
- route withdrawal/drain occurs before removal of active traffic;
- regional/local preference is advisory and never overrides health policy;
- measured path quality is preferred over static assumptions;
- retries are bounded and idempotent where applicable;
- every privileged routing mutation remains subject to identity, policy, approval and audit.

## Deployment

The service fabric is designed to integrate existing FTN_service and redy capabilities rather than duplicate them. A new service node is registered, validated, health-checked and made eligible without reconstructing `familytimenet.com` DNS or the global mesh.
