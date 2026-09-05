# Contributing to FTN-AI

Contributions should preserve the production contract of FTN-AI: changes must be reviewable, testable, and safe to deploy.

## Before opening a PR

1. Branch from `main`.
2. Make the smallest change that fixes the underlying problem.
3. Run `gofmt` for changed Go files.
4. Run the affected Go tests; for control-plane changes run `go test ./...` from `services/control-plane`.
5. Validate shell scripts with `bash -n` and Compose files with `docker compose config --quiet` when applicable.
6. Do not commit credentials, private keys, generated secrets, or production connection strings.

## Security and operations

- Privileged network, router, OLT, deployment, and destructive operations must remain approval-gated.
- CI may validate, test, build, and collect evidence; it must not perform real infrastructure mutations.
- Preserve tenant isolation, identity boundaries, audit trails, idempotency, and durable execution invariants.
- Avoid weakening authentication, authorization, TLS, or validation checks to make a test pass.

## Pull requests

Describe the root cause, the production impact, the files changed, and the validation performed. If a check could not be executed, state that explicitly rather than claiming success.

## Reporting bugs

Include the failing command or endpoint, expected behavior, actual behavior, relevant logs, and a minimal reproduction when possible. Never include passwords, access tokens, private keys, or other secrets.
