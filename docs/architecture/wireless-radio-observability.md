# FTN Wireless + Radio Observability

FTN treats wireless/radio sensing as an optional edge telemetry plane. It is not part of the default control-plane startup.

## Sensor classes

| Class | FTN role | Primary collector |
|---|---|---|
| Wi-Fi 2.4/5/6 GHz | authorized WLAN telemetry/WIDS | Kismet Linux Wi-Fi datasource |
| Bluetooth/BLE | authorized local device discovery telemetry | Kismet Linux HCI / supported BLE capture datasource |
| Zigbee | authorized IoT-radio telemetry | Kismet supported Zigbee datasource |
| RTL-SDR / sub-GHz sensors | authorized sensor telemetry | Kismet `rtl_433` datasource |
| SDR-capable Wi-Fi hardware | specialized wireless capture | Kismet supported SDR datasource |

Kismet 2025-09-R1 supports multiple datasource families, including Linux Wi-Fi, Linux Bluetooth HCI, Zigbee adapters, RTL-SDR/rtl_433 and other SDR-backed sources. The exact datasource must match the installed hardware and capture tool; FTN does not invent a source type.

## Deployment boundary

- The `wireless` Compose profile is opt-in.
- Kismet uses host networking because wireless capture interfaces and local radio devices are host resources.
- Capture privileges are limited to the sensor boundary; the control plane receives approved telemetry rather than arbitrary packet data.
- No raw wireless/radio collection is enabled by default.
- Only FTN-owned/managed and explicitly authorized sensor environments may be configured.
- Customer-facing services do not receive direct access to sensor hardware.

## Source configuration

`deployment/observability/kismet_ftn_sources.conf` contains the supported source patterns. Enable only the lines matching hardware actually installed on the sensor node.

Do not copy a source type between hardware families. For example, RTL-SDR/rtl_433 requires an SDR device; a normal Wi-Fi or Bluetooth adapter cannot substitute for it.

## Operational layers

```text
radio / Wi-Fi / BLE / Zigbee hardware
            |
            v
       Kismet sensor
            |
            v
 approved telemetry adapter
            |
            +----> FTN event / metrics pipeline
            |
            +----> local retention / audit
            v
     FTN control plane
```

The wireless/radio plane is therefore an observability input, not a privileged route-changing mechanism. Any resulting network mutation still crosses the normal authorization, approval, verification and audit boundary.
