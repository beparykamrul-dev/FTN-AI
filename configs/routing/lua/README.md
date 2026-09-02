# FTN Lua Routing Policy

Lua is an embedded policy language for deterministic, reviewable routing decisions.

## Rules

- Lua never executes shell commands.
- Lua cannot access the network, filesystem, environment, or process APIs.
- Lua receives normalized route/flow metadata only.
- A Lua result is advisory until the FTN approval workflow authorizes a route change.
- CPU/time/resource limits must be enforced by the host runtime.
- Every policy execution is auditable by policy version and decision hash.

## Intended inputs

`prefix`, `asn`, `country`, `latency_class`, `traffic_class`, `peer_role`, `route_source`.

## Intended outputs

`allow`, `deny`, `local_preference`, `med`, `community`, `reason`.

The production host should use a sandboxed Lua 5.1-compatible runtime and treat policy files as configuration, not executable deployment scripts.
