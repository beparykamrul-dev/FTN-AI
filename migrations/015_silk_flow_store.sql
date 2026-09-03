BEGIN;

CREATE TABLE IF NOT EXISTS ftn_silk_flow_records (
    id BIGSERIAL PRIMARY KEY,
    observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    tenant_id TEXT NOT NULL,
    exporter_id TEXT NOT NULL,
    version SMALLINT NOT NULL,
    source_ip INET,
    destination_ip INET,
    source_port INTEGER,
    destination_port INTEGER,
    protocol INTEGER,
    bytes BIGINT NOT NULL DEFAULT 0,
    packets BIGINT NOT NULL DEFAULT 0,
    sampling_rate INTEGER NOT NULL DEFAULT 1,
    fingerprint TEXT NOT NULL,
    service_id TEXT,
    subscriber_id TEXT,
    main_server_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ftn_silk_flow_version_ck CHECK (version IN (5, 9, 10)),
    CONSTRAINT ftn_silk_flow_sampling_ck CHECK (sampling_rate > 0),
    CONSTRAINT ftn_silk_flow_bytes_ck CHECK (bytes >= 0),
    CONSTRAINT ftn_silk_flow_packets_ck CHECK (packets >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS ftn_silk_flow_dedupe_idx
    ON ftn_silk_flow_records (tenant_id, exporter_id, fingerprint, observed_at);

CREATE INDEX IF NOT EXISTS ftn_silk_flow_time_idx
    ON ftn_silk_flow_records (tenant_id, observed_at DESC);

CREATE INDEX IF NOT EXISTS ftn_silk_flow_src_dst_idx
    ON ftn_silk_flow_records (source_ip, destination_ip, observed_at DESC);

CREATE INDEX IF NOT EXISTS ftn_silk_flow_service_idx
    ON ftn_silk_flow_records (service_id, main_server_id, observed_at DESC);

COMMIT;
