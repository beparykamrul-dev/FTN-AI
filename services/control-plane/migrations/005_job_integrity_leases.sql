-- FTN durable job integrity and lease hardening. Additive/idempotent.
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS request_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS lease_expires_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS durable_jobs_lease_idx ON durable_jobs(status, lease_expires_at);
CREATE INDEX IF NOT EXISTS durable_jobs_idempotency_tenant_idx ON durable_jobs(tenant_id, idempotency_key);

-- Preserve the existing global idempotency constraint while recording the exact
-- request identity needed to reject accidental key reuse with different data.
