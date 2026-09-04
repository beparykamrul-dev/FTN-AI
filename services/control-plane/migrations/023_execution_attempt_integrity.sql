-- Enforce durable execution integrity at the database transaction boundary.
-- The deferred constraint trigger allows a claim/finish transaction to update
-- durable_jobs and execution_attempts in either order, but the transaction
-- cannot commit a running job without its matching running attempt.

CREATE OR REPLACE FUNCTION ftn_assert_running_job_attempt() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.status = 'running' THEN
    IF NOT EXISTS (
      SELECT 1
        FROM execution_attempts ea
       WHERE ea.job_id = NEW.id
         AND ea.attempt_no = NEW.attempts
         AND ea.worker_id = COALESCE(NEW.locked_by, '')
         AND ea.status = 'running'
    ) THEN
      RAISE EXCEPTION 'running durable job % has no matching execution attempt %', NEW.id, NEW.attempts
        USING ERRCODE = '23514';
    END IF;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_execution_integrity ON durable_jobs;
CREATE CONSTRAINT TRIGGER durable_jobs_execution_integrity
AFTER INSERT OR UPDATE OF status, attempts, locked_by
ON durable_jobs
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION ftn_assert_running_job_attempt();

INSERT INTO schema_migrations(version,name)
VALUES (23,'023_execution_attempt_integrity.sql')
ON CONFLICT (version) DO NOTHING;
