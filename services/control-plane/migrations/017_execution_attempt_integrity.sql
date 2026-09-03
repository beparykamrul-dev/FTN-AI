-- FTN execution-attempt integrity boundary. Additive/idempotent.
CREATE OR REPLACE FUNCTION ftn_validate_execution_attempt()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE job_status TEXT; job_worker TEXT;
BEGIN
  IF NEW.job_id IS NULL OR NULLIF(BTRIM(NEW.worker_id),'') IS NULL THEN RAISE EXCEPTION 'execution_attempt_binding_required'; END IF;
  IF NEW.attempt_no <= 0 THEN RAISE EXCEPTION 'execution_attempt_number_invalid'; END IF;
  SELECT status,COALESCE(locked_by,'') INTO job_status,job_worker FROM durable_jobs WHERE id=NEW.job_id FOR SHARE;
  IF NOT FOUND THEN RAISE EXCEPTION 'durable_job_not_found'; END IF;
  IF TG_OP='INSERT' THEN
    IF job_status<>'running' THEN RAISE EXCEPTION 'execution_attempt_requires_running_job'; END IF;
    IF job_worker IS DISTINCT FROM NEW.worker_id THEN RAISE EXCEPTION 'execution_attempt_worker_mismatch'; END IF;
  ELSE
    IF NEW.job_id IS DISTINCT FROM OLD.job_id OR NEW.attempt_no IS DISTINCT FROM OLD.attempt_no OR NEW.worker_id IS DISTINCT FROM OLD.worker_id THEN RAISE EXCEPTION 'execution_attempt_identity_immutable'; END IF;
    IF OLD.status IN ('succeeded','failed') AND NEW.status IS DISTINCT FROM OLD.status THEN RAISE EXCEPTION 'finished_execution_attempt_immutable'; END IF;
    IF NEW.status='running' AND job_status<>'running' THEN RAISE EXCEPTION 'running_attempt_requires_running_job'; END IF;
  END IF;
  RETURN NEW;
END; $$;
DROP TRIGGER IF EXISTS execution_attempt_integrity ON execution_attempts;
CREATE TRIGGER execution_attempt_integrity BEFORE INSERT OR UPDATE OF job_id,attempt_no,worker_id,status ON execution_attempts FOR EACH ROW EXECUTE FUNCTION ftn_validate_execution_attempt();
CREATE INDEX IF NOT EXISTS execution_attempts_worker_idx ON execution_attempts(worker_id,started_at DESC);
