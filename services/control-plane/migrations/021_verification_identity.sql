-- Persist the authenticated verifier for durable execution auditability.
ALTER TABLE durable_jobs
  ADD COLUMN IF NOT EXISTS verified_by UUID REFERENCES principals(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS durable_jobs_verified_by_idx
  ON durable_jobs(tenant_id, verified_by);

INSERT INTO schema_migrations(version,name)
VALUES (21,'021_verification_identity.sql')
ON CONFLICT (version) DO NOTHING;
