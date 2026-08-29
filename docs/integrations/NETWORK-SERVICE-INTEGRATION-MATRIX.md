# FTN Network & Service Integration Matrix

This document records provider-neutral integration boundaries for network, service-discovery, traffic-analysis, and connectivity components. External projects are integration references/dependencies, not copied source code.

## Network fabric

- Kube-OVN: optional cloud-native SDN/CNI adapter. Use for Kubernetes/KubeVirt VPCs, subnets, VLAN/underlay, IPAM, NAT/EIP, QoS, distributed gateways, BGP exposure and traffic mirroring. It is not a transport protocol. Kube-OVN supports multi-tenant VPC networking and BGP exposure. 
- SD-Core: optional mobile-core integration boundary; keep its control/data-plane responsibilities separate from FTN's generic control plane.

## Service fabric

- PolarisMesh: optional service discovery/governance adapter. Integrate service registration, health checks, routing/load-balancing policy, rate limiting, circuit breaking and configuration governance through FTN's Service Registry. Do not duplicate a second authoritative service registry without an explicit federation policy.
- Polaris Console capabilities should remain behind FTN Control Panel rather than becoming a second administrator surface.

## Traffic intelligence

- SiLK/NetFlow/IPFIX tooling: optional flow-ingestion and historical traffic-analysis adapter. Store normalized flow metadata in FTN's traffic-intelligence boundary; do not couple the control plane directly to a collector implementation.

## Connectivity

- curl/libcurl: generic HTTP/HTTPS/API transfer capability for adapters, health probes and controlled external integrations. It is not a service-placement mechanism.

## Future transport registry

The following remain provider-neutral transport adapters and are intentionally not forced into the core build until the platform foundation is stable:

- WireGuard
- AmneziaWG
- OpenVPN 3
- GRE
- SSLH
- Shadowsocks
- Hysteria2
- TUIC
- Socket/WebSocket
- BATMAN-adv
- Yggdrasil
- CJDNS

## Placement policy

Every integration is represented as a capability. The FTN Control Plane selects eligible servers based on health, CPU/RAM/storage capacity, network capacity, locality, policy, redundancy and explicit authorization. Service placement is not hard-coded to a single server.

High-impact infrastructure operations remain approval-gated and auditable. This matrix does not grant automatic deployment permission.
