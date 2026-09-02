-- FTN durable job lifecycle journal. Additive/idempotent.
-- Database trigger keeps job state and its lifecycle event in the same transaction.

CREATE OR REPLACE FUNCTION ftn_journal_job_transition() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  next_sequence BIGINT;
  transition TEXT;
BEGIN
  IF TG_OP = 'INSERT' THEN
    transition := 'job.submitted';
  ELSIF NEW.status IS DISTINCT FROM OLD.status THEN
    transition := CASE NEW.status
      WHEN 'running' THEN 'job.claimed'
      WHEN 'succeeded' THEN 'job.succeeded'
      WHEN 'failed' THEN 'job.failed'
      WHEN 'queued' THEN 'job.requeued'
      WHEN 'cancelled' THEN 'job.cancelled'
      ELSE 'job.state_changed'
    END;
  ELSE
    RETURN NEW;
  END IF;

  PERFORM pg_advisory_xact_lock(hashtext(COALESCE(NEW.tenant_id::text, '')));
  SELECT COALESCE(MAX(sequence) + 1, 1) INTO next_sequence
    FROM event_journal
   WHERE tenant_id IS NOT DISTINCT FROM NEW.tenant_id;

  INSERT INTO event_journal(
    tenant_id,event_type,sequence,correlation_id,causation_id,aggregate_id,payload
  ) VALUES (
    NEW.tenant_id,
    transition,
    next_sequence,
    NEW.correlation_id,
    '',
    NEW.id::text,
    jsonb_build_object(
      'job_id', NEW.id::text,
      'job_type', NEW.job_type,
      'status', NEW.status,
      'attempts', NEW.attempts,
      'max_attempts', NEW.max_attempts,
      'execution_action', COALESCE(NEW.execution_action, ''),
      'approval_id', COALESCE(NEW.approval_id::text, ''),
      'locked_by', COALESCE(NEW.locked_by, ''),
      'last_error', COALESCE(NEW.last_error, '')
    )
  );
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_event_journal ON durable_jobs;
CREATE TRIGGER durable_jobs_event_journal
AFTER INSERT OR UPDATE OF status ON durable_jobs
FOR EACH ROW EXECUTE FUNCTION ftn_journal_job_transition();
