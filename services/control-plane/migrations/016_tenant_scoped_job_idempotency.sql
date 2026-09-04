-- Make durable-job idempotency storage tenant-scoped without changing the API contract.
-- The application supplies the logical key; the database stores a tenant-prefixed key.
-- This preserves the existing global UNIQUE constraint while preventing a same-key
-- request from one tenant from replaying a job owned by another tenant.

UPDATE durable_jobs
SET idempotency_key = tenant_id::text || ':' || idempotency_key
WHERE tenant_id IS NOT NULL
  AND left(idempotency_key, length(tenant_id::text) + 1) <> tenant_id::text || ':';

CREATE OR REPLACE FUNCTION ftn_scope_durable_job_idempotency_key() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  prefix TEXT;
BEGIN
  IF NEW.tenant_id IS NULL THEN
    RAISE EXCEPTION 'durable-job tenant_id is required';
  END IF;
  prefix := NEW.tenant_id::text || ':';
  IF left(NEW.idempotency_key, length(prefix)) <> prefix THEN
    NEW.idempotency_key := prefix || NEW.idempotency_key;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_scope_idempotency ON durable_jobs;
CREATE TRIGGER durable_jobs_scope_idempotency
BEFORE INSERT OR UPDATE OF tenant_id, idempotency_key
ON durable_jobs
FOR EACH ROW EXECUTE FUNCTION ftn_scope_durable_job_idempotency_key();

INSERT INTO schema_migrations(version,name)
VALUES (16,'016_tenant_scoped_job_idempotency.sql')
ON CONFLICT (version) DO NOTHING;
