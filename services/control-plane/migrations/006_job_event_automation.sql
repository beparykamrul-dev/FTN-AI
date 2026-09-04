-- FTN durable job lifecycle journal. Additive/idempotent.
-- Canonical database triggers keep job/approval state and lifecycle events in the same transaction.

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

  PERFORM pg_advisory_xact_lock(hashtextextended(COALESCE(NEW.tenant_id::text, ''), 0));
  SELECT COALESCE(MAX(sequence) + 1, 1) INTO next_sequence
    FROM event_journal
   WHERE tenant_id IS NOT DISTINCT FROM NEW.tenant_id;

  INSERT INTO event_journal(
    tenant_id,event_type,sequence,correlation_id,causation_id,aggregate_id,payload
  ) VALUES (
    NEW.tenant_id,
    transition,
    next_sequence,
    COALESCE(NEW.correlation_id, ''),
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

DROP TRIGGER IF EXISTS durable_jobs_lifecycle_event ON durable_jobs;
DROP TRIGGER IF EXISTS durable_jobs_event_trigger ON durable_jobs;
DROP TRIGGER IF EXISTS durable_jobs_event_journal ON durable_jobs;

CREATE TRIGGER durable_jobs_event_journal
AFTER INSERT OR UPDATE OF status ON durable_jobs
FOR EACH ROW EXECUTE FUNCTION ftn_journal_job_transition();

CREATE OR REPLACE FUNCTION ftn_approval_lifecycle_event() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  event_name TEXT;
BEGIN
  IF NEW.status IS NOT DISTINCT FROM OLD.status THEN
    RETURN NEW;
  END IF;

  event_name := CASE NEW.status
    WHEN 'approved' THEN 'approval.approved'
    WHEN 'rejected' THEN 'approval.rejected'
    WHEN 'expired' THEN 'approval.expired'
    WHEN 'executed' THEN 'approval.executed'
    WHEN 'rolled_back' THEN 'approval.rolled_back'
    ELSE 'approval.state.changed'
  END;

  PERFORM pg_advisory_xact_lock(hashtextextended(COALESCE(NEW.tenant_id::text, ''), 0));

  INSERT INTO event_journal(
    tenant_id, event_type, sequence, correlation_id, causation_id,
    aggregate_id, payload
  )
  VALUES (
    NEW.tenant_id,
    event_name,
    COALESCE((SELECT MAX(sequence) + 1 FROM event_journal WHERE tenant_id = NEW.tenant_id), 1),
    '',
    '',
    NEW.id::text,
    jsonb_build_object(
      'approval_id', NEW.id,
      'action', NEW.action,
      'resource', NEW.resource,
      'status', NEW.status,
      'requested_by', NEW.requested_by,
      'approved_by', NEW.approved_by
    )
  );

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_lifecycle_event ON change_approvals;
CREATE TRIGGER change_approvals_lifecycle_event
AFTER UPDATE OF status
ON change_approvals
FOR EACH ROW EXECUTE FUNCTION ftn_approval_lifecycle_event();

INSERT INTO schema_migrations(version,name)
VALUES (6,'006_job_event_automation.sql')
ON CONFLICT (version) DO NOTHING;
