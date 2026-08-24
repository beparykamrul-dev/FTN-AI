# FTN ElastiFlow NetObservability

## Purpose

Integrate ElastiFlow NetObserv capabilities as a provider/reference adapter for FTN network-flow observability.

## Ingest

- NetFlow
- IPFIX
- sFlow
- AWS VPC Flow Logs where applicable
- Device/application telemetry

## FTN pipeline

Flow exporters -> FTN Flow Gateway -> ElastiFlow/NetObserv adapter -> normalize/enrich -> OTel/event fabric -> ClickHouse/Elasticsearch/OpenSearch -> FTN AIOps

## Outputs

- Elasticsearch / Elastic Cloud
- OpenSearch
- Kafka / compatible event streams
- Splunk / Cribl-compatible adapters where required
- FTN canonical flow schema

## Correlation

Correlate flow records with router, OLT, BGP, DNS, Geo/S2, latency, health, provider, POP, application and customer-service metadata.

## Production principles

- Provider-neutral adapter boundary
- No mandatory third-party dependency in FTN core
- License/provenance tracked per component
- mTLS for protected data paths
- Sampling and retention policies configurable
- Raw payloads are not required for ordinary flow analytics
- High-volume flow data is externalized from application memory
- Backpressure, batching and bounded queues
- Health checks and automatic recovery

## Related ElastiFlow references

The current ElastiFlow NetObserv Helm chart describes a Unified Flow Collector that decodes, transforms, normalizes, translates and enriches IPFIX, NetFlow, sFlow and AWS VPC Flow Logs, with outputs including Elasticsearch, OpenSearch and Kafka-family platforms.

The legacy Logstash-based ElastiFlow repository is archived; FTN should reference the current generation rather than treating the legacy repository as the production dependency.

## FTN AIOps signals

- top talkers
- source/destination conversations
- application/service traffic
- ASN traffic
- exporter health
- interface/POP traffic
- traffic locality
- capacity trends
- anomaly signals
- latency/packet-loss correlation

## Architecture

Go data-plane services and ASP.NET Core control-plane services consume the canonical FTN flow model. AIOps consumes normalized flow events and may recommend or execute approved remediation workflows according to FTN policy.
