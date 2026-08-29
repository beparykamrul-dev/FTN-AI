# FTN Service Placement

The control plane treats a server as a capability-bearing node rather than a fixed home for a service.

## Endpoints

- `GET /api/v1/nodes` — current node capabilities and health.
- `POST /api/v1/nodes/register` — register/update a node capability snapshot.
- `POST /api/v1/placement/preview` — rank eligible nodes for a service.

Placement considers service capability, health, CPU, RAM, SSD, HDD, network throughput, latency, packet loss, region and provider affinity. The result is a **plan only**; production execution remains approval-gated.

## Example node registration

```json
{"id":"pop-bd-01","provider":"provider-a","region":"BD","services":["ftndns","media","proxy"],"cpu_percent":18,"ram_percent":32,"ssd_percent":20,"hdd_percent":40,"net_mbps":9000,"latency_ms":4,"packet_loss_percent":0.1,"healthy":true}
```

The same registry is intended to be consumed by DNS, proxy, tunnel, media, API, communication and future transport adapters so service placement does not depend on a hard-coded server.
