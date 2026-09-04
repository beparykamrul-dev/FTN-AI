-- Make approval request deduplication tenant-scoped.
-- The application already includes tenant/principal identity in request_hash;
-- the database constraint must enforce the same isolation boundary.
ALTER TABLE change_approvals
  DROP CONSTRAINT IF EXISTS change_approvals_request_hash_key;

CREATE UNIQUE INDEX IF NOT EXISTS change_approvals_tenant_request_hash_uq
  ON change_approvals(tenant_id, request_hash);

INSERT INTO schema_migrations(version,name)
VALUES (14,'014_tenant_scoped_approval_hash.sql')
ON CONFLICT (version) DO NOTHING;
