-- Make governed data requests replay-safe without altering legacy rows.
ALTER TABLE data_requests
  ADD COLUMN IF NOT EXISTS request_hash TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS data_requests_tenant_request_hash_uidx
  ON data_requests(tenant_id, request_hash)
  WHERE request_hash IS NOT NULL;

INSERT INTO schema_migrations(version,name)
VALUES (22,'022_data_request_idempotency.sql')
ON CONFLICT (version) DO NOTHING;
