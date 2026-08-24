# FTN Security & Network Observability Fabric

## Purpose

Unify network security telemetry, flow monitoring, route health, investigation notebooks, and operational analytics behind provider-neutral FTN contracts.

## Components

- **Zeek**: network security telemetry and protocol metadata boundary.
- **FastNetMon**: high-speed DDoS/anomaly detection and mitigation signal boundary.
- **SIEM**: normalized security-event ingestion, correlation, retention, and alerting boundary.
- **RoutePulse**: route/path-health adapter boundary for BGP and reachability intelligence; implementation remains replaceable.
- **Jupyter**: controlled investigation and analytics workspace, consuming sanitized/authorized datasets rather than production secrets.
- **Recoil**: not a core FTN runtime dependency. If this refers to the archived Facebook Recoil state-management project, it is excluded from the backend/network core; modern FTN UI state should use an actively maintained abstraction.

## Flow

`Packet/Flow -> Zeek/FastNetMon/eBPF -> Event Normalizer -> SIEM/ClickHouse -> AIOps -> Alert/Route/Capacity Decision`

`BGP/Route telemetry -> Route adapter -> Path score -> Geo/latency engine -> service placement/routing`

`Sanitized telemetry -> Jupyter -> investigation/report -> approved operational change`

## Design rules

- Security tools remain isolated from the customer data plane.
- No single vendor becomes an FTN core dependency.
- Raw telemetry and derived metrics have separate retention policies.
- Secrets, private keys, credentials, and personal data are excluded from notebooks and exported telemetry.
- Mitigation actions pass through the existing FTN policy/control boundary.
- FastNetMon/Zeek outputs are normalized into FTN-native event contracts.
- SIEM is an integration boundary, not the authoritative FTN database.
- ClickHouse remains the high-volume analytics path; SIEM/search stores serve security investigation use cases.
- Route health is correlated with flow, DNS, Geo/S2, RTT, loss, and POP health.
- Jupyter is for controlled investigation, reproducibility, and offline analytics—not unrestricted production execution.

## Operational loop

`Observe -> Correlate -> Detect -> Score -> Recommend -> Policy Gate -> Mitigate/Route/Scale -> Verify -> Record`
