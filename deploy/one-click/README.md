# FTN One-Click Server Bootstrap

Run from the repository root:

```bash
sudo bash deploy/one-click/bootstrap.sh
```

The bootstrap prepares a Debian/Ubuntu-style host, installs required base dependencies when available, validates Docker Compose configuration, creates runtime directories, and starts the Compose stack when a compose manifest exists.

## Extension model

New FTN services should be added as independently versioned Compose services/modules. The bootstrap intentionally does not hard-code provider credentials or service-specific secrets.

## Production requirements

Before public exposure, configure TLS, DNS, firewall policy, backups, monitoring, secret management, and provider credentials through the deployment environment/secret store. Validate the stack with `docker compose config` before rollout.
