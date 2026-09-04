CREATE DATABASE IF NOT EXISTS ftn_analytics;

CREATE TABLE IF NOT EXISTS ftn_analytics.flow_records
(
    observed_at DateTime64(3, 'UTC'),
    tenant_id String,
    exporter_id String,
    version UInt16,
    source_ip String,
    destination_ip String,
    source_port UInt16,
    destination_port UInt16,
    protocol UInt8,
    bytes UInt64,
    packets UInt64,
    sampling_rate UInt32,
    fingerprint String,
    service_id String,
    subscriber_id String,
    main_server_id String
)
ENGINE = MergeTree
PARTITION BY toDate(observed_at)
ORDER BY (tenant_id, exporter_id, observed_at, fingerprint)
TTL observed_at + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;
