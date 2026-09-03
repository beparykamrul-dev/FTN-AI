-- FTN execution-attempt replay/idempotency boundary. Additive/idempotent.
-- A terminal attempt is a one-way event. A second terminalization is rejected.

CREATE OR REPLACE FUNCTION ftn_reject_execution_attempt_replay()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF OLD.status IN ('succeeded','failed') THEN
        IF NEW.status IS DISTINCT FROM OLD.status
           OR NEW.job_id IS DISTINCT FROM OLD.job_id
           OR NEW.attempt_no IS DISTINCT FROM OLD.attempt_no
           OR NEW.worker_id IS DISTINCT FROM OLD.worker_id THEN
            RAISE EXCEPTION 'terminal_execution_attempt_immutable';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempt_terminal_immutability ON execution_attempts;
CREATE TRIGGER execution_attempt_terminal_immutability
BEFORE UPDATE OF status, job_id, attempt_no, worker_id
ON execution_attempts
FOR EACH ROW
EXECUTE FUNCTION ftn_reject_execution_attempt_replay();

-- A job can have at most one attempt row for a given attempt number.
-- This closes duplicate-attempt insertion/replay for the same durable execution.
CREATE UNIQUE INDEX IF NOT EXISTS execution_attempts_job_attempt_unique
    ON execution_attempts(job_id, attempt_no);

-- A worker may not create a second active attempt for the same job.
CREATE UNIQUE INDEX IF NOT EXISTS execution_attempts_active_job_unique
    ON execution_attempts(job_id)
 WHERE status = 'running';

-- A terminal attempt must remain tied to its original job/attempt identity.
CREATE OR REPLACE FUNCTION ftn_validate_attempt_identity()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.attempt_no <= 0 THEN
        RAISE EXCEPTION 'execution_attempt_number_invalid';
    END IF;
    IF NULLIF(BTRIM(NEW.worker_id), '') IS NULL THEN
        RAISE EXCEPTION 'execution_attempt_worker_required';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempt_identity_integrity ON execution_attempts;
CREATE TRIGGER execution_attempt_identity_integrity
BEFORE INSERT OR UPDATE OF job_id, attempt_no, worker_id
ON execution_attempts
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_attempt_identity();
