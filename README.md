# FTN-AI

FTN-AI is the private, policy-controlled intelligence and automation layer for the Family Time Network (FTN) platform.

## Production goals

- Lightweight agent runtime with service-, user- and developer-scoped agents.
- Approval-first execution for privileged or destructive actions.
- Private/local runtime support with provider-neutral model boundaries.
- Unified control-plane integration for DNS, network operations, fiber/GIS, monitoring, mesh, proxy, mail, SMS and FTN services.
- Versioned contracts as the source of truth for platform behavior.
- Web and Android client integration through stable APIs.

## Repository layout

- `internal/platform/` — reusable FTN platform engines and control-plane primitives.
- `services/` — deployable FTN service boundaries.
- `modules/` — bounded business modules.
- `backend/` — diagnostics and supporting services.
- `frontend/` — control-center contracts and UI architecture.
- `web/` — web and mobile client surfaces.
- `os/ftn-os/` — FTN runtime/module lifecycle layer.
- `configs/v1/` — versioned runtime and policy contracts.
- `contracts/v1/` — control-center API/data contracts.
- `docs/architecture/` — platform architecture and completion specifications.
- `scripts/` — repository and server-agent validation utilities.

## Safety and control boundary

Agent responses may request an action, but privileged execution must pass through the approval/control boundary. Secrets remain server-side and telemetry must not export secrets or private user data. CI is validation-only for real infrastructure: it must not mutate routers, OLTs, or production network state.

## Validation

The repository gates validate shell syntax, YAML/JSON, Go formatting, control-plane tests/builds, Compose configuration, security boundaries, API authentication, durable execution, recovery behavior, and architecture contracts. Go CI uses the repository's Go 1.24 toolchain where applicable.

Run the affected checks locally before opening a PR. If a check cannot be executed, report that explicitly.

## Development principle

This repository targets real FTN infrastructure. Avoid demo-only implementations, fake production telemetry, embedded credentials, silent privileged execution, or destructive automation without an explicit policy/approval path.
