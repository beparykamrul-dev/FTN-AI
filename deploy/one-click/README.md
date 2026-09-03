# FTN One-Click Production Deployment

FTN uses one canonical deployment path for the repository's production-marked Compose stacks.

## Fresh Debian/Ubuntu server

For a specific branch:

```bash
curl -fsSL https://raw.githubusercontent.com/beparykamrul-dev/FTN-AI/fix/migration-validation-015-023/install.sh | sudo env FTN_REF=fix/migration-validation-015-023 bash
```

After this work is merged to `main`, use:

```bash
curl -fsSL https://raw.githubusercontent.com/beparykamrul-dev/FTN-AI/main/install.sh | sudo bash
```

## Existing checkout

```bash
sudo bash deploy/one-click/bootstrap.sh
```

The bootstrap prepares the host and hands off to `deploy/one-click/live.sh`. The live runner discovers every `docker-compose.yml` / `compose.yml` explicitly marked with:

```yaml
x-ftn-production-stack: true
```

It validates every discovered stack before changing runtime state, starts the control-plane/migration foundation first, then starts the remaining production stacks, waits for control-plane readiness, and prints final service status.

## Runtime guarantees

- Docker Engine and Compose v2 are required.
- Runtime secrets are generated locally when missing and stored in `.env` with mode `0600`.
- Existing secrets are preserved; credentials are never printed.
- Production discovery is opt-in; test/development Compose files are not started accidentally.
- PostgreSQL data uses persistent Compose volumes.
- Database migrations run through the declared migration runner.
- Privileged FTN operations remain policy/approval-gated.
- Deployment fails closed if no production stack is marked or a Compose validation/readiness check fails.

## Adding another production service

Put its Compose manifest anywhere in the repository and add:

```yaml
x-ftn-production-stack: true
```

Then redeploy with:

```bash
sudo bash deploy/one-click/live.sh
```

No second deployment script is required.

## Public exposure

One-click deployment does not silently change DNS, router ACLs, firewall policy, certificates, or provider credentials. Those changes remain explicit and environment-specific. TLS/reverse proxy, DNS, firewall, backups, monitoring, and external-provider credentials must be configured through their approved FTN deployment components.
