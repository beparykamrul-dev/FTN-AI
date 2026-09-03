# FTN One-Click Production Bootstrap

The installer is idempotent: it installs host prerequisites, Docker Engine/Compose, clones or updates FTN-AI, generates missing runtime secrets locally, validates Compose, builds the control plane, starts PostgreSQL + migrations + control-plane, and waits for readiness.

## Fresh Debian/Ubuntu server

After this branch is merged to `main`:

```bash
curl -fsSL https://raw.githubusercontent.com/beparykamrul-dev/FTN-AI/main/install.sh | sudo bash
```

For branch testing:

```bash
curl -fsSL https://raw.githubusercontent.com/beparykamrul-dev/FTN-AI/repair/ci-build-2026-09-03/install.sh | sudo env FTN_REF=repair/ci-build-2026-09-03 bash
```

## Existing checkout

```bash
sudo bash deploy/one-click/bootstrap.sh
```

Default installation directory: `/opt/ftn-ai`. Override with `FTN_INSTALL_DIR=/srv/ftn-ai`.

## Runtime behavior

- Docker is enabled and started through systemd.
- PostgreSQL data persists in the named Compose volume.
- Missing `FTN_DB_PASSWORD` and `FTN_API_AUTH_TOKEN` are generated locally with OpenSSL and stored only in `.env` (`0600`).
- Compose is validated before startup.
- Bootstrap waits for `/readyz` or `/healthz` and prints service diagnostics if readiness fails.
- Credentials are never printed and provider credentials are not embedded in the repository.
- Privileged FTN operations remain policy/approval-gated.

## Public production exposure

The one-click bootstrap starts the FTN control plane but does not silently change DNS, router ACLs, firewall policy, or provider credentials. Public exposure should add the approved TLS/reverse-proxy, DNS, firewall, backup, monitoring, and secret-management configuration for the target FTN environment.
