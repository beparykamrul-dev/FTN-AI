# FTN Observability Services

This stack adds FTN-owned observability building blocks:

- **Kismet** — authorized wireless telemetry/sensor input for FTN-managed networks. It is isolated behind the `wireless` Compose profile and is not part of the default control-plane startup.
- **Kismet radio datasources** — opt-in Wi-Fi, Bluetooth/BLE, Zigbee and RTL-SDR/`rtl_433` sensor inputs where the installed hardware and capture tool support them. Source patterns live in `kismet_ftn_sources.conf`.
- **Kopia** — backup repository/server for FTN backup and disaster-recovery workflows. It uses the `backup` profile and persistent volumes.
- **OpenSearch + Dashboards** — searchable log/event/observability backend. It uses the `observability` profile and binds its host ports to loopback by default.

## Profiles

```text
observability  -> OpenSearch + OpenSearch Dashboards
backup         -> Kopia
wireless       -> Kismet + authorized wireless/radio sensor boundary
```

The core FTN control plane remains independent of these services. They can therefore fail without becoming a single point of failure for the control API.

## Wireless/radio scope

The radio plane is an **observability input**, not a route-changing authority. Supported sensor classes include:

- Wi-Fi 2.4/5/6 GHz through Kismet Linux Wi-Fi capture
- Bluetooth/BLE discovery or supported capture hardware
- Zigbee adapters supported by Kismet
- RTL-SDR/`rtl_433` sensor telemetry
- Other Kismet-supported SDR-backed sources when the required hardware/capture component is actually installed

Only explicitly authorized FTN-owned/managed sensor environments may enable these sources. Raw packet/radio collection is not enabled by the baseline configuration.

See `docs/architecture/wireless-radio-observability.md` for the boundary and `deployment/observability/kismet_ftn_sources.conf` for source configuration patterns.

## Data flow

```text
FTN services / approved telemetry
          |
          +--> OpenSearch (logs/events/search)
          |
          +--> Kopia (backup/DR repository)
          |
          +--> Kismet (authorized wireless/radio telemetry)
```

## Security boundary

- Secrets are environment-driven; never commit real passwords.
- OpenSearch and Kopia are bound to `127.0.0.1` on the host in this baseline and should be exposed externally only through the FTN reverse-proxy/authentication layer.
- Kismet is a privileged sensor component. Keep it limited to FTN-owned/managed radio interfaces and use privilege separation appropriate to the host.
- Do not collect or retain data outside the authorized FTN monitoring scope.
- Production backup repositories must be encrypted and restore-tested.
- Hardware-specific capture components must be matched to their documented Kismet datasource; FTN does not assume that one radio adapter can perform every capture type.

## Version policy

OpenSearch is pinned to `3.8.0`. Kopia is configurable through `KOPIA_VERSION`; production deployments must replace `latest` with a reviewed immutable release tag.

Kismet is pinned to stable release `2025-09-R1` in the baseline configuration.
