-- Enforce the inverse side of execution-attempt integrity.
-- Any durable job committed in running state must have exactly one matching
-- running attempt with the same attempt number and worker.

CREATE OR REPLACE FUNCTION ftn_assert_running_job_matches_attempt() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  attempt_count INTEGER;
BEGIN
  IF NEW.status = 'running' THEN
    SELECT count(*)
      INTO attempt_count
      FROM execution_attempts ea
     WHERE ea.job_id = NEW.id
       AND ea.attempt_no = NEW.attempts
       AND ea.worker_id = COALESCE(NEW.locked_by, '')
       AND ea.status = 'running';

    IF attempt_count <> 1 THEN
      RAISE EXCEPTION 'running durable job % does not have exactly one matching running execution attempt',
        NEW.id USING ERRCODE = '23514';
    END IF;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_job_running_attempt_binding ON durable_jobs;
CREATE CONSTRAINT TRIGGER durable_job_running_attempt_binding
AFTER INSERT OR UPDATE OF status, attempts, locked_by
ON durable_jobs
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION ftn_assert_running_job_matches_attempt();

INSERT INTO schema_migrations(version,name)
VALUES (26,'026_durable_job_running_attempt_binding.sql')
ON CONFLICT (version) DO NOTHING;
