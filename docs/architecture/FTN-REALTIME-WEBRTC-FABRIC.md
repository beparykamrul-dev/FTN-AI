# FTN Realtime WebRTC Fabric

## Components

- LiveKit: primary Go-based realtime SFU/reference layer.
- mediasoup: optional Node.js/C++ SFU adapter/reference for workloads needing its media-routing model.
- coturn: STUN/TURN connectivity layer for NAT traversal and relay fallback.

## FTN architecture

Client -> FTN Realtime Gateway -> Geo/Latency Steering -> SFU -> STUN/TURN -> Provider/Client

### Capabilities

- WebRTC audio/video/data
- UDP-first transport with TCP/TURN fallback
- STUN candidate discovery
- TURN relay for restricted NAT/firewall paths
- SFU-based selective forwarding
- simulcast/SVC-aware media policies
- congestion and packet-loss telemetry
- latency/jitter/bitrate monitoring
- provider/POP-aware realtime routing
- autoscaling and health-aware SFU placement
- OpenTelemetry/Prometheus telemetry
- ClickHouse analytics integration
- FTN AI voice/realtime assistant integration
- recording/egress adapter boundary
- Android/PC/Web/TV client compatibility

## Production policy

LiveKit, mediasoup, and coturn remain replaceable provider-neutral components behind FTN interfaces. FTN owns routing, identity, policy, telemetry, health, automation, and control-plane integration.

Coturn provides standard STUN/TURN relay functionality and supports UDP/TCP/IPv6/TLS-related WebRTC connectivity. LiveKit provides a scalable Go SFU and has built-in TURN/network fallback capabilities. mediasoup remains an alternative SFU implementation with Node.js server integration and C++ media core.
