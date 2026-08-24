# FTN Next-Gen Realtime Edge Fabric

## Purpose
Provider-neutral realtime, edge, media, telemetry and autonomous operations layer for FTN.

## Core capabilities
- WebRTC/realtime media adapter boundary
- HTTP/1.1, HTTP/2, HTTP/3, QUIC, TCP and UDP
- Go data-plane and ASP.NET Core control-plane
- REST, gRPC, WebSocket, SSE and webhook APIs
- mTLS/PKI and short-lived service credentials
- EDA, asynchronous jobs, queues, backpressure and durable workflows
- Edge compute, serverless and workload placement
- Geo-steering + S2Geometry + latency/health scoring
- BGP/DNS/path intelligence
- eBPF/XDP telemetry and traffic controls
- Prometheus + OpenTelemetry + ClickHouse/ClickStack telemetry path
- ELK/Elasticsearch adapter for log/search workloads
- Redis for bounded ephemeral cache/state only
- AIOps: anomaly detection, correlation, RCA, prediction and recovery orchestration
- Autoscaling based on traffic, latency, queue depth and resource health
- Self-healing: detect -> isolate -> reroute/replace -> verify
- Android, PC, Web and TV client integration boundaries
- Provider adapters kept isolated from FTN core

## Realtime/media boundary
Daily-style realtime/WebRTC capability is represented as a provider adapter, not as a hard dependency. FTN can substitute another realtime provider or self-hosted media implementation without changing the control plane.

## Design rule
Open-source capability first; official provider credentials/API integrations remain optional adapters. Persistent state is externalized; request/session processing is stateless-first.
