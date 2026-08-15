# FTN Source Consolidation Registry

This registry records source material discovered from the user's existing FTN planning/source files before consolidation.

## Confirmed source areas

- Sprint-01 foundation: SOURCE, DOCS, TESTS, DOCKER, DATABASE, PREVIEW.
- Backend/API foundation and AI chat foundation.
- Prometheus, SNMP Exporter, NetFlow/nfcapd, LibreNMS, OTDR importer, SmokePing, Syslog-ng.
- Leaflet/GeoJSON fiber and customer topology mapping.
- Grafana integration and multi-datasource monitoring.
- Fiber → OLT → ONU → Switch → Customer topology.
- PPPoE heatmap/live users.
- Billing and customer management.
- AI anomaly/predictive alert concepts.
- Network/fiber/customer maps and NetBox integration.
- Multi-panel ISP management: Admin, Employee, Reseller, User.
- AI Assistant, AI Automation and AI Call Center concepts.

## Consolidation rules

1. Preserve source material; do not silently delete competing implementations.
2. Prefer production implementations over placeholders/prototypes when source evidence supports that choice.
3. Keep legacy material under an explicit archive/reference area when it may still contain useful functionality.
4. Track duplicates, missing dependencies and unresolved conflicts before merging.
5. Keep FTN-owned services and integrations modular so external providers remain optional adapters.
6. Privileged autonomous network changes remain policy/authorization controlled.

## Known source evidence

The existing source material documents SNMP, NetFlow, LibreNMS, OTDR, SmokePing, Syslog-ng, Leaflet, Prometheus/Grafana and NetBox integration, plus multi-panel ISP management and AI capabilities.
