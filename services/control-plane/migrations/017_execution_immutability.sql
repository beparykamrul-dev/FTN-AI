-- FTN worker execution immutability boundary. Additive/idempotent.
-- Once a durable job is approval-bound, the worker cannot replace the reviewed
-- execution object while changing any approval-bound field.
CREATE OR REPLACE FUNCTION ftn_enforce_approved_job_immutability()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.approval_id IS NOT NULL THEN
        IF NEW.approval_id IS DISTINCT FROM OLD.approval_id
           OR NEW.execution_action IS DISTINCT FROM OLD.execution_action
           OR NEW.approval_request_hash IS DISTINCT FROM OLD.approval_request_hash
           OR NEW.execution_payload_hash IS DISTINCT FROM OLD.execution_payload_hash
           OR NEW.payload::jsonb IS DISTINCT FROM OLD.payload::jsonb THEN
            RAISE EXCEPTION 'approved_execution_object_immutable';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_execution_immutability ON durable_jobs;
CREATE TRIGGER durable_jobs_execution_immutability
BEFORE UPDATE OF approval_id, execution_action, approval_request_hash, execution_payload_hash, payload
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_enforce_approved_job_immutability();

-- Fail closed if an approved job somehow reaches execution without its
-- persisted approval binding. This is intentionally enforced at the database
-- boundary as well as in the application claim path.
CREATE OR REPLACE FUNCTION ftn_enforce_approved_job_claim()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    approval_status TEXT;
    approval_expires TIMESTAMPTZ;
    approval_hash TEXT;
    approval_payload JSONB;
BEGIN
    IF NEW.approval_id IS NULL OR NEW.status <> 'running' OR OLD.status = 'running' THEN
        RETURN NEW;
    END IF;

    SELECT status, expires_at, request_hash, approval_payload
      INTO approval_status, approval_expires, approval_hash, approval_payload
      FROM change_approvals
     WHERE id = NEW.approval_id
       AND tenant_id IS NOT DISTINCT FROM NEW.tenant_id
     FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'approval_not_found_or_tenant_mismatch';
    END IF;
    IF approval_status IS DISTINCT FROM 'approved'
       OR (approval_expires IS NOT NULL AND approval_expires <= now()) THEN
        RAISE EXCEPTION 'approval_not_approved_or_expired';
    END IF;
    IF NEW.approval_request_hash IS DISTINCT FROM approval_hash THEN
        RAISE EXCEPTION 'approval_request_hash_mismatch';
    END IF;
    IF NEW.payload::jsonb IS DISTINCT FROM approval_payload THEN
        RAISE EXCEPTION 'execution_payload_approval_payload_mismatch';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_approved_claim_gate ON durable_jobs;
CREATE TRIGGER durable_jobs_approved_claim_gate
BEFORE UPDATE OF status ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_enforce_approved_job_claim();
