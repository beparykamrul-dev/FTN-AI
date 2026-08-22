# FTN Network Kernel v1

FTN's Jupyter-facing network kernel is a controlled wrapper around approved network tooling. It is designed for diagnostics, inventory, telemetry collection and explicitly authorized changes.

## Integration points

- `ipykernel.kernelbase.Kernel`: notebook/kernel lifecycle boundary.
- `do_execute`: mapped to FTN's structured `KernelToolRegistry`, not arbitrary Python or shell forwarding.
- Netmiko: device CLI collection adapter.
- Paramiko: SSH transport adapter where required by an approved integration.
- NAPALM: normalized multi-vendor network state/collection adapter.
- NetSA/CERT tooling: flow/measurement research and telemetry integration point.

## Execution model

`Jupyter request -> FTN kernel wrapper -> capability/policy -> registered tool -> device/provider -> observed result -> audit`

The wrapper deliberately does not expose a generic `exec()` or shell escape. Device credentials remain outside notebooks and are supplied through the server-side secret boundary. Collection operations are read-only by default. Configuration-changing operations require the FTN approval policy and post-change verification.

## Tool naming

The registry uses stable provider-neutral names (`netmiko`, `paramiko`, `napalm`, `cert-netsa`, `ipykernel`). Actual Python packages are deployment dependencies and are not embedded into the Go control plane.

## Compatibility

The kernel wrapper can coexist with the existing FTN control plane, DNS, mesh and telemetry modules. It is an adapter layer, not a second orchestration system.
