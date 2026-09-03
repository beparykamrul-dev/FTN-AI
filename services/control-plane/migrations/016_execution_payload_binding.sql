-- FTN approved execution payload binding. Additive/idempotent.
ALTER TABLE change_approvals ADD COLUMN IF NOT EXISTS approval_payload_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS execution_payload_hash TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS durable_jobs_execution_payload_hash_idx ON durable_jobs(execution_payload_hash);
CREATE INDEX IF NOT EXISTS change_approvals_payload_hash_idx ON change_approvals(approval_payload_hash);

-- Backfill legacy rows and make the payload hash deterministic at the database boundary.
UPDATE change_approvals
   SET approval_payload_hash = md5(approval_payload::text)
 WHERE NULLIF(BTRIM(approval_payload_hash), '') IS NULL;

CREATE OR REPLACE FUNCTION ftn_bind_execution_payload_to_approval()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    approval_payload JSONB;
    approval_payload_hash TEXT;
BEGIN
    IF NEW.approval_id IS NULL THEN
        NEW.execution_payload_hash := '';
        RETURN NEW;
    END IF;

    SELECT a.approval_payload, a.approval_payload_hash
      INTO approval_payload, approval_payload_hash
      FROM change_approvals a
     WHERE a.id = NEW.approval_id
       AND a.tenant_id IS NOT DISTINCT FROM NEW.tenant_id
     FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'approval_not_found_or_tenant_mismatch';
    END IF;
    IF NEW.payload::jsonb IS DISTINCT FROM approval_payload THEN
        RAISE EXCEPTION 'execution_payload_approval_payload_mismatch';
    END IF;
    IF NULLIF(BTRIM(approval_payload_hash), '') IS NULL THEN
        RAISE EXCEPTION 'approval_payload_hash_missing';
    END IF;

    -- The hash is copied only after exact JSONB equality succeeds. This prevents
    -- a caller from supplying a different hash for a different execution object.
    NEW.execution_payload_hash := approval_payload_hash;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_execution_payload_binding ON durable_jobs;
CREATE TRIGGER durable_jobs_execution_payload_binding
BEFORE INSERT OR UPDATE OF approval_id, tenant_id, payload, execution_payload_hash
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_bind_execution_payload_to_approval();

CREATE OR REPLACE FUNCTION ftn_set_approval_payload_hash()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.approval_payload IS NULL THEN
        RAISE EXCEPTION 'approval_payload_required';
    END IF;
    NEW.approval_payload_hash := md5(NEW.approval_payload::text);
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_payload_hash ON change_approvals;
CREATE TRIGGER change_approvals_payload_hash
BEFORE INSERT OR UPDATE OF approval_payload
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_set_approval_payload_hash();
