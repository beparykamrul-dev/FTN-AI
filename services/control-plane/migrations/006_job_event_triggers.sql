-- FTN durable job lifecycle -> Event Journal bridge.
-- Additive/idempotent. Keeps lifecycle events in the same DB transaction as the mutation.

CREATE OR REPLACE FUNCTION ftn_job_lifecycle_event() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  event_name TEXT;
  tenant UUID;
BEGIN
  tenant := NEW.tenant_id;

  IF TG_OP = 'INSERT' THEN
    event_name := 'job.submitted';
  ELSIF NEW.status IS DISTINCT FROM OLD.status THEN
    event_name := CASE NEW.status
      WHEN 'running' THEN 'job.claimed'
      WHEN 'succeeded' THEN 'job.finished'
      WHEN 'failed' THEN 'job.finished'
      WHEN 'cancelled' THEN 'job.cancelled'
      WHEN 'queued' THEN CASE
        WHEN NEW.execution_action LIKE '%.rollback' THEN 'job.rollback.requested'
        ELSE 'job.requeued'
      END
      ELSE 'job.state.changed'
    END;
  ELSIF NEW.verification_payload IS DISTINCT FROM OLD.verification_payload
        AND NEW.verification_payload <> '{}'::jsonb THEN
    event_name := 'job.verified';
  ELSE
    RETURN NEW;
  END IF;

  PERFORM pg_advisory_xact_lock(hashtext(COALESCE(tenant::text, '')));

  INSERT INTO event_journal(
    tenant_id, event_type, sequence, correlation_id, causation_id,
    aggregate_id, payload
  )
  VALUES (
    tenant,
    event_name,
    COALESCE((SELECT MAX(sequence) + 1 FROM event_journal WHERE tenant_id = tenant), 1),
    COALESCE(NEW.correlation_id, ''),
    '',
    NEW.id::text,
    jsonb_build_object(
      'job_id', NEW.id,
      'job_type', NEW.job_type,
      'status', NEW.status,
      'attempts', NEW.attempts,
      'max_attempts', NEW.max_attempts,
      'approval_id', NEW.approval_id,
      'execution_action', NEW.execution_action,
      'worker_id', NEW.locked_by,
      'last_error', NEW.last_error
    )
  );

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_lifecycle_event ON durable_jobs;
CREATE TRIGGER durable_jobs_lifecycle_event
AFTER INSERT OR UPDATE OF status, verification_payload
ON durable_jobs
FOR EACH ROW EXECUTE FUNCTION ftn_job_lifecycle_event();

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

  PERFORM pg_advisory_xact_lock(hashtext(COALESCE(NEW.tenant_id::text, '')));

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
