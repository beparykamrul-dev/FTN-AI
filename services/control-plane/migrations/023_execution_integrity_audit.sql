-- FTN execution integrity audit boundary. Additive/idempotent.
-- Makes the approval/claim trigger ordering explicit so the shared advisory
-- lock is acquired before the older approval validation trigger executes.

DROP TRIGGER IF EXISTS durable_jobs_approval_execution_lock ON durable_jobs;
CREATE TRIGGER durable_jobs_approval_execution_lock
BEFORE UPDATE OF status
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_lock_approval_for_running_job();

-- Keep terminal execution objects immutable after they are finalized.
CREATE OR REPLACE FUNCTION ftn_freeze_terminal_job_execution_metadata()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.status IN ('succeeded','failed','cancelled') THEN
        IF NEW.payload::jsonb IS DISTINCT FROM OLD.payload::jsonb
           OR NEW.approval_id IS DISTINCT FROM OLD.approval_id
           OR NEW.execution_action IS DISTINCT FROM OLD.execution_action
           OR NEW.execution_payload_hash IS DISTINCT FROM OLD.execution_payload_hash
           OR NEW.approval_request_hash IS DISTINCT FROM OLD.approval_request_hash THEN
            RAISE EXCEPTION 'terminal_execution_metadata_immutable';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_terminal_metadata_immutability ON durable_jobs;
CREATE TRIGGER durable_jobs_terminal_metadata_immutability
BEFORE UPDATE OF payload, approval_id, execution_action, execution_payload_hash, approval_request_hash
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_freeze_terminal_job_execution_metadata();

-- A successful privileged job must retain a successful execution attempt.
-- This is checked at the database boundary in addition to the application gate.
CREATE OR REPLACE FUNCTION ftn_require_successful_attempt_for_approved_job()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    attempt_status TEXT;
BEGIN
    IF NEW.status = 'succeeded' AND NEW.approval_id IS NOT NULL AND OLD.status IS DISTINCT FROM NEW.status THEN
        SELECT status INTO attempt_status
          FROM execution_attempts
         WHERE job_id = NEW.id
           AND attempt_no = NEW.attempts
         LIMIT 1;
        IF attempt_status IS DISTINCT FROM 'succeeded' THEN
            RAISE EXCEPTION 'approved_job_requires_successful_execution_attempt';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_success_attempt_gate ON durable_jobs;
CREATE TRIGGER durable_jobs_success_attempt_gate
BEFORE UPDATE OF status
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_require_successful_attempt_for_approved_job();
