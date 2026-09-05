-- Keep execution_attempts state bound to the durable job state.
-- This closes the inverse integrity gap left by the deferred durable_jobs trigger:
-- an attempt cannot be marked/routed as running unless its job is running
-- under the same worker and attempt number.

CREATE OR REPLACE FUNCTION ftn_assert_running_attempt_matches_job() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  job_status TEXT;
  job_attempts INTEGER;
  job_worker TEXT;
BEGIN
  IF NEW.status = 'running' THEN
    SELECT status, attempts, COALESCE(locked_by, '')
      INTO job_status, job_attempts, job_worker
      FROM durable_jobs
     WHERE id = NEW.job_id;

    IF NOT FOUND
       OR job_status <> 'running'
       OR job_attempts <> NEW.attempt_no
       OR job_worker <> COALESCE(NEW.worker_id, '') THEN
      RAISE EXCEPTION 'running execution attempt %/% does not match durable job state',
        NEW.job_id, NEW.attempt_no USING ERRCODE = '23514';
    END IF;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS execution_attempts_job_integrity ON execution_attempts;
CREATE CONSTRAINT TRIGGER execution_attempts_job_integrity
AFTER INSERT OR UPDATE OF status, attempt_no, worker_id
ON execution_attempts
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION ftn_assert_running_attempt_matches_job();

INSERT INTO schema_migrations(version,name)
VALUES (25,'025_execution_attempt_state_integrity.sql')
ON CONFLICT (version) DO NOTHING;
