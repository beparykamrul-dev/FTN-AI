-- Bind approvals to the exact request payload and retire the duplicate
-- durable-job lifecycle trigger installed by the legacy 006 migration.
ALTER TABLE change_approvals
  ADD COLUMN IF NOT EXISTS payload_hash TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS change_approvals_payload_hash_idx
  ON change_approvals(tenant_id, payload_hash);

DROP TRIGGER IF EXISTS durable_jobs_lifecycle_event ON durable_jobs;
