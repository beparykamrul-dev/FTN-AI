# FTN Mesh

The mesh package owns the normalized topology and approval boundary for FTN mesh transports.

## Adapter status

BATMAN-adv, Yggdrasil, and CJDNS are **not registered as executable host adapters yet**. They remain capability declarations until their platform-specific management APIs, privilege model, rollback semantics, and integration tests are implemented.

This is intentional: a capability registry entry must not be interpreted as proof that a live transport driver exists.

## Required adapter boundary

A concrete adapter must implement `MeshAdapter` and provide:

- read-only topology/health discovery;
- deterministic input validation;
- immutable plan identity;
- approval-bound mutation;
- pre-change state capture;
- post-change verification;
- rollback of the exact approved plan;
- bounded command/API execution with context cancellation;
- structured errors and audit events.

Adapters must never accept arbitrary shell commands, credentials, or an unapproved route mutation through the mesh package.

## Host integration rule

The FTN control plane may discover that a host supports a mesh capability, but support discovery is not admission. An adapter becomes executable only after its implementation and integration tests are present and the capability is explicitly admitted by the FTN registry/policy layer.
