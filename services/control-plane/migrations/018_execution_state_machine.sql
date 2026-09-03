-- FTN durable execution state-machine boundary. Additive/idempotent.
-- A job may only move through the states defined by the durable worker lifecycle.

CREATE OR REPLACE FUNCTION ftn_validate_durable_job_transition()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.status IS NOT DISTINCT FROM OLD.status THEN
        RETURN NEW;
    END IF;

    IF OLD.status = 'queued' AND NEW.status IN ('running','cancelled') THEN
        RETURN NEW;
    END IF;
    IF OLD.status = 'running' AND NEW.status IN ('queued','succeeded','failed','cancelled') THEN
        RETURN NEW;
    END IF;

    RAISE EXCEPTION 'invalid_durable_job_transition:%->%', OLD.status, NEW.status;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_state_machine ON durable_jobs;
CREATE TRIGGER durable_jobs_state_machine
BEFORE UPDATE OF status
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_durable_job_transition();

-- A worker can finish only the attempt that owns the current job lease.
-- The terminal job state must agree with the attempt state.
CREATE OR REPLACE FUNCTION ftn_validate_durable_job_terminal_state()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    attempt_status TEXT;
    attempt_worker TEXT;
BEGIN
    IF NEW.status NOT IN ('succeeded','failed') OR OLD.status = NEW.status THEN
        RETURN NEW;
    END IF;

    SELECT ea.status, ea.worker_id
      INTO attempt_status, attempt_worker
      FROM execution_attempts ea
     WHERE ea.job_id = NEW.id
       AND ea.attempt_no = NEW.attempts
     ORDER BY ea.id DESC
     LIMIT 1;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'terminal_job_requires_matching_execution_attempt';
    END IF;
    IF attempt_status IS DISTINCT FROM NEW.status THEN
        RAISE EXCEPTION 'job_attempt_terminal_state_mismatch';
    END IF;
    IF NULLIF(BTRIM(NEW.locked_by), '') IS NOT NULL
       AND attempt_worker IS DISTINCT FROM NEW.locked_by THEN
        RAISE EXCEPTION 'terminal_job_worker_mismatch';
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_terminal_integrity ON durable_jobs;
CREATE TRIGGER durable_jobs_terminal_integrity
BEFORE UPDATE OF status
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_durable_job_terminal_state();

-- Terminal attempts are immutable; running attempts may only become terminal.
CREATE OR REPLACE FUNCTION ftn_validate_execution_attempt_transition()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.status IS NOT DISTINCT FROM OLD.status THEN
        RETURN NEW;
    END IF;

    IF OLD.status = 'running' AND NEW.status IN ('succeeded','failed') THEN
        RETURN NEW;
    END IF;

    RAISE EXCEPTION 'invalid_execution_attempt_transition:%->%', OLD.status, NEW.status;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempt_state_machine ON execution_attempts;
CREATE TRIGGER execution_attempt_state_machine
BEFORE UPDATE OF status
ON execution_attempts
FOR EACH ROW
EXECUTE FUNCTION ftn_validate_execution_attempt_transition();

-- Finished attempts must carry a completion timestamp.
CREATE OR REPLACE FUNCTION ftn_require_execution_attempt_finished_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.status IN ('succeeded','failed') AND NEW.finished_at IS NULL THEN
        RAISE EXCEPTION 'finished_execution_attempt_requires_finished_at';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempt_finished_at_integrity ON execution_attempts;
CREATE TRIGGER execution_attempt_finished_at_integrity
BEFORE INSERT OR UPDATE OF status, finished_at
ON execution_attempts
FOR EACH ROW
EXECUTE FUNCTION ftn_require_execution_attempt_finished_at();
