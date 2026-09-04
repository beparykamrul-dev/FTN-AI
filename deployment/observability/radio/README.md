# FTN Radio / Wireless Sensor Plane

This directory defines the production boundary for authorized FTN wireless and radio telemetry.

## Sources

- Kismet: managed Wi-Fi/BLE/802.15.4 wireless telemetry.
- RTL-SDR / rtl_433: optional external sensor adapters for compatible, authorized sensor telemetry.
- GPS: optional sensor-location metadata.

## Architecture

```text
Managed radio interfaces / external sensors
                 |
        source-specific adapter
                 |
              Kismet
                 |
        FTN telemetry adapter
          /       |        \
     metrics    events      audit
       |           |           |
  Prometheus   OpenSearch  event-journal
                 |
            FTN control plane
```

The sensor plane is observational by default. It does not provide arbitrary RF transmission, unauthorized collection, or a bypass around FTN approval controls.

## Deployment

`wireless` and `sdr` are opt-in profiles. They must never be silently enabled on a production node merely because the observability stack is enabled.

Hardware-specific datasource configuration belongs on the target sensor host. Do not put interface names, frequencies, device serials, or credentials into the repository.

## Data policy

- Collect only telemetry required for FTN operations and authorized monitoring.
- Prefer metadata over raw captures.
- Raw IQ/radio capture is disabled by default and requires an explicit authorization boundary.
- Sensor configuration mutations remain approval-bound.
