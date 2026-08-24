# FTN Reseller Service Policy

## Purpose

FTN services are modular. A reseller may sell one service, multiple services, or a bundled service according to the reseller's authorized commercial and network policy.

## Service modes

- `NETWORK_ONLY`: FTN network-layer service without requiring FTN-provided customer data storage, OLT, fiber, or wireless access infrastructure.
- `ROUTER_ONLY`: authorized router/gateway service.
- `OLT_ACCESS`: authorized OLT/PON service.
- `FIBER_ACCESS`: authorized wired fiber access service.
- `WIRELESS_ACCESS`: authorized wireless access service.
- `DATA_STORAGE`: authorized customer data/storage service.
- `BUNDLE`: a policy-defined combination of the above.

## Reseller policy

A reseller profile defines which service modes it may offer. The control plane must enforce the profile rather than assuming one universal delivery method.

A reseller may be both:

1. a consumer of an FTN service; and
2. an authorized reseller of that or another FTN service.

Bandwidth allocations may also be resold when the reseller has an explicit bandwidth-resale entitlement. Parent allocation, child allocation, actual usage, burst policy, and remaining capacity must be reconciled before a child offer is activated.

## Network-only semantics

`NETWORK_ONLY` does not imply a physical access medium. The reseller supplies or arranges its authorized access path according to its own policy. FTN provides the network-layer service and enforces identity, authorization, routing policy, service limits, health, billing state, and audit requirements.

No OLT, fiber, wireless link, or FTN customer-data storage is required by the `NETWORK_ONLY` service contract itself.

An actual Internet/service path still requires an authorized upstream or interconnection and appropriate network transport; the service contract must not claim physical connectivity where none exists.

## Hierarchy

`FTN -> Master Reseller -> Reseller/Sub-reseller -> Customer`

Each child relationship records its parent, service entitlement, bandwidth entitlement, policy, billing relationship, and authorization state.

## Control-plane requirements

The control plane must validate:

- reseller authorization;
- service-mode entitlement;
- parent/child relationship;
- bandwidth allocation and remaining capacity;
- routing eligibility;
- node/service health;
- billing status;
- audit requirements.

A requested service must not be provisioned when the reseller policy does not permit that service mode or allocation.

## Safety boundary

Provider peering, BGP sessions, credentials, certificates, customer access, and physical infrastructure remain subject to explicit authorization and applicable provider/regulatory requirements. The policy layer does not bypass those controls.
