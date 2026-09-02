# FTN Observability Services

This stack adds three FTN-owned building blocks:

- **Kismet** — authorized wireless telemetry/sensor input for FTN-managed networks. It is isolated behind the `wireless` Compose profile and is not part of the default control-plane startup.
- **Kopia** — backup repository/server for FTN backup and disaster-recovery workflows. It uses the `backup` profile and persistent volumes.
- **OpenSearch + Dashboards** — searchable log/event/observability backend. It uses the `observability` profile and binds its host ports to loopback by default.

## Profiles

```text
observability  -> OpenSearch + OpenSearch Dashboards
backup         -> Kopia
wireless       -> Kismet
```

The core FTN control plane remains independent of these services. They can therefore fail without becoming a single point of failure for the control API.

## Data flow

```text
FTN services / approved telemetry
          |
          +--> OpenSearch (logs/events/search)
          |
          +--> Kopia (backup/DR repository)
          |
          +--> Kismet (authorized wireless telemetry)
```

## Security boundary

- Secrets are environment-driven; never commit real passwords.
- OpenSearch and Kopia are bound to `127.0.0.1` on the host in this baseline and should be exposed externally only through the FTN reverse-proxy/authentication layer.
- Kismet is a privileged sensor component. Keep it limited to FTN-owned/managed radio interfaces and use privilege separation appropriate to the host.
- Do not collect or retain data outside the authorized FTN monitoring scope.
- Production backup repositories must be encrypted and restore-tested.

## Version policy

OpenSearch is pinned to `3.8.0`. Kopia is configurable through `KOPIA_VERSION`; production deployments must replace `latest` with a reviewed immutable release tag.

Kismet is pinned to the current stable release `2025-09-R1` in the baseline configuration.
