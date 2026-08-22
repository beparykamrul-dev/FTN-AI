# FTN-AI

FTN-AI is the private, policy-controlled intelligence and automation layer for the Family Time Network (FTN) platform.

## Production goals

- One lightweight agent runtime with service-scoped, user-scoped and developer-scoped agents.
- Approval-first execution: conversational AI does not silently perform privileged or destructive actions.
- Private/local runtime support with provider-neutral model boundaries.
- Unified control-plane integration for DNS, network operations, fiber/GIS, monitoring, mesh, proxy, mail, SMS and FTN services.
- Versioned YAML contracts and architecture documents as the source of truth for platform behavior.
- Web and Android client integration through stable contracts.

## Repository layout

- `internal/platform/` — reusable FTN platform engines and control-plane primitives.
- `services/` — deployable FTN service boundaries.
- `modules/` — bounded business modules.
- `backend/` — backend diagnostics and supporting services.
- `frontend/` — control-center contracts and UI architecture.
- `web/` — web and mobile client surfaces.
- `os/ftn-os/` — FTN runtime/module lifecycle layer.
- `configs/v1/` — versioned runtime and policy contracts.
- `contracts/v1/` — control-center API/data contracts.
- `docs/architecture/` — platform architecture and completion specifications.
- `scripts/` — source organization and server-agent utilities.

## Safety and control boundary

The FTN AI layer is designed around explicit policy enforcement. Agent responses may request an action, but privileged execution must pass through the existing approval/control boundary. Secrets must remain server-side and telemetry must not export secrets or private user data.

## Validation

The repository CI validates repository structure, YAML syntax, shell syntax and the Go module currently declared under `modules/account`. Additional Go source trees are intentionally kept as independently bounded platform/service packages until their module boundaries are declared.

## Development principle

This repository is intended for real FTN infrastructure. Avoid demo-only implementations, fake production telemetry, embedded credentials, silent privileged execution, or destructive automation without an explicit policy/approval path.
