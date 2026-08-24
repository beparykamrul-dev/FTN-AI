# FTN Database Fabric

The database layer is workload-aware and provider-neutral.

- PostgreSQL: transactional/control data
- ClickHouse: telemetry/analytics
- Redis: ephemeral cache/session workloads
- Object storage: files, media and backups
- Local/POP databases: latency-sensitive edge workloads

A database broker evaluates compatibility, health, latency, capacity, quota, replication and policy before placement. Free capacity alone never authorizes migration or paid-resource usage.
