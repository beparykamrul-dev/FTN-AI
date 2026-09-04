BEGIN;

ALTER TABLE ftn_silk_flow_records
    ADD COLUMN IF NOT EXISTS first_seen TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_seen TIMESTAMPTZ;

ALTER TABLE ftn_silk_flow_records
    DROP CONSTRAINT IF EXISTS ftn_silk_flow_time_range_ck;

ALTER TABLE ftn_silk_flow_records
    ADD CONSTRAINT ftn_silk_flow_time_range_ck
    CHECK (first_seen IS NULL OR last_seen IS NULL OR last_seen >= first_seen);

CREATE INDEX IF NOT EXISTS ftn_silk_flow_first_last_idx
    ON ftn_silk_flow_records (tenant_id, first_seen DESC, last_seen DESC);

COMMIT;
