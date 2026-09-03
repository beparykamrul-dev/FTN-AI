-- FTN approval execution race boundary. Additive/idempotent.
-- Prevents a queued privileged job from being claimed after its approval is
-- revoked/expired, and prevents approval mutation once execution has begun.

CREATE OR REPLACE FUNCTION ftn_lock_approval_for_running_job()
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
     FOR UPDATE;

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

DROP TRIGGER IF EXISTS durable_jobs_approval_execution_lock ON durable_jobs;
CREATE TRIGGER durable_jobs_approval_execution_lock
BEFORE UPDATE OF status ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_lock_approval_for_running_job();

-- Once a privileged job has started, its approval cannot be revoked or edited
-- through the approval state machine. This is intentionally fail-closed.
CREATE OR REPLACE FUNCTION ftn_freeze_approval_while_job_running()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    active_jobs INTEGER;
BEGIN
    SELECT COUNT(*)
      INTO active_jobs
      FROM durable_jobs j
     WHERE j.approval_id = OLD.id
       AND j.status = 'running';

    IF active_jobs > 0 THEN
        IF NEW.status IS DISTINCT FROM OLD.status
           OR NEW.action IS DISTINCT FROM OLD.action
           OR NEW.resource IS DISTINCT FROM OLD.resource
           OR NEW.request_hash IS DISTINCT FROM OLD.request_hash
           OR NEW.approval_payload IS DISTINCT FROM OLD.approval_payload THEN
            RAISE EXCEPTION 'approval_locked_by_running_execution';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_running_execution_lock ON change_approvals;
CREATE TRIGGER change_approvals_running_execution_lock
BEFORE UPDATE OF status, action, resource, request_hash, approval_payload
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_freeze_approval_while_job_running();

CREATE INDEX IF NOT EXISTS durable_jobs_approval_status_idx
    ON durable_jobs(approval_id, status);
