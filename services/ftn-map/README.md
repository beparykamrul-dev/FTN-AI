# FTN Map — Fiber Intelligence

Native FTN fiber topology and recovery layer.

## Topology objects

- fiber paths
- tables/cabinets
- joint boxes
- splitters
- OLT/ONU relationships
- routers
- PPPoE user profiles
- cut/fault events
- distance calculations

The map layer is designed for PostgreSQL/PostGIS and can expose topology through the FTN API.

## AI fiber intelligence

AI analyzes topology, distance, path risk and fault signals to produce advisory cut detection and recovery recommendations.

Recovery actions remain policy-gated and authorized; AI cannot independently make privileged network changes.

## Local recovery

The design supports local-first topology/cache state so a POP or client node can continue to display known topology and recovery context during upstream loss, then reconcile with the FTN control plane when connectivity returns.
