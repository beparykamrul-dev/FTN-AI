# FTN Packet & Edge Fabric

## Purpose

This document defines the provider-neutral integration boundary for low-level packet, edge, overlay, telemetry, and protocol-validation capabilities in FTN-AI.

These capabilities are **adapters and governed services**, not a replacement for the FTN control plane. Production changes remain subject to identity, authorization, policy validation, approval, execution, verification, and audit.

## Layer model

1. **Hardware / edge**
   - CAN_RAW / PF_CAN and can-utils for authorized industrial/edge nodes.
   - Device capability discovery and health telemetry.
2. **Kernel packet path**
   - eBPF/XDP for policy enforcement, telemetry and fast-path classification.
   - SOCK_RAW only for explicitly authorized packet diagnostics, protocol research and controlled adapters.
   - No raw-packet feature bypasses FTN authorization or host/network policy.
3. **Overlay / SDN**
   - Kube-OVN / OVS.
   - VXLAN.
   - IPsec tunnels.
   - Placement and route policy remain owned by the control plane.
4. **Transport adapters**
   - WireGuard, AmneziaWG, OpenVPN 3, GRE, SSLH, Shadowsocks, Hysteria2, Socket/WebSocket, BATMAN-adv, Yggdrasil and CJDNS are registered as transport capabilities and can be expanded independently after the core platform is stable.
5. **Flow telemetry**
   - NetFlow v5/v9, IPFIX, YAF and SiLK/rwflowpack-compatible pipelines.
   - Flow data is normalized into FTN telemetry contracts.
6. **Protocol validation / security**
   - Parsing-discrepancy detection is a defensive validation function.
   - Test harnesses compare parser interpretations across approved gateway/WAF/proxy/backend implementations.
   - The system records ambiguity, rejects unsafe normalization, and creates hardening findings.
   - It must not provide production bypass instructions or deploy bypass payloads.

## Control-plane integration

Every adapter exposes:

- capability name and version
- supported platforms
- required privileges/capabilities
- health/readiness state
- resource requirements
- telemetry outputs
- configuration schema
- rollback requirements
- security policy requirements

The control panel selects **service + eligible server/node + deployment policy**. Placement is calculated from health, capacity, locality, latency, policy and redundancy rather than from a hard-coded host list.

## Resource placement

The scheduler must consider CPU, RAM, SSD/HDD capacity, network capacity, node health and service constraints. A service can run on one or many eligible nodes. Migration requires a governed change set and verification.

## Protocol discrepancy testing

The FTN security validation plane should include parser differential tests for protocols actually exposed by FTN services. The test result model should capture:

- protocol and parser versions
- input class identifier
- parser A normalized view
- parser B normalized view
- discrepancy classification
- severity
- affected boundary
- recommended hardening action
- regression-test identifier

Sensitive test cases remain isolated from production traffic and credentials.

## Platform references

The implementation may use existing open-source projects where their licenses and deployment requirements permit. Provider/project names are metadata; FTN does not imply ownership, endorsement or bundled licensing.

Examples include Linux PF_CAN, can-utils, Kube-OVN, OVS, VXLAN/IPsec, NetFlow/IPFIX tooling, YAF and SiLK. Perl/Python, OpenRC, Cygwin and CDN libraries may be used where they are appropriate deployment dependencies.

## Acceptance criteria

A feature is not marked production-ready merely because its source or adapter exists. It requires configuration validation, tests, security/policy tests, health/metrics contracts, failure handling, audit behavior, idempotent registration and documented deployment/rollback behavior.
