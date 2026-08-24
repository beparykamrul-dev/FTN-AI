# FTN Observability Fabric

Metrics, logs, traces and network-flow telemetry are treated as separate but correlated signals. Collection is local-first at POP/edge nodes and forwarded according to retention and capacity policy.

`Device -> Collector -> Normalizer -> Buffer -> Metrics/Logs/Traces/Flow -> Correlation -> NOC/SOC/AIOps`

Core requirements: health checks, cardinality control, local buffering, retention policies, tenant isolation, alert deduplication and auditability.
