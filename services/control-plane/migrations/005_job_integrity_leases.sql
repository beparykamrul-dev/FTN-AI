-- FTN durable job integrity and lease hardening. Additive/idempotent migration.
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS request_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS lease_expires_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS durable_jobs_lease_idx ON durable_jobs(status, lease_expires_at);
CREATE INDEX IF NOT EXISTS durable_jobs_idempotency_tenant_idx ON durable_jobs(tenant_id, idempotency_key);

-- Preserve the existing global idempotency constraint while recording the exact
-- request identity needed to reject accidental key reuse with different data.

-- Atomic durable job lifecycle events. State changes and their replay events
-- commit together because this trigger executes in the same database transaction.
CREATE OR REPLACE FUNCTION ftn_job_event_trigger() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  event_name TEXT;
  event_payload JSONB;
  tenant TEXT;
  corr TEXT;
BEGIN
  tenant := COALESCE(NEW.tenant_id::text, '');
  corr := COALESCE(NEW.correlation_id, '');

  IF TG_OP = 'INSERT' THEN
    event_name := 'job.submitted';
  ELSIF OLD.status IS DISTINCT FROM NEW.status THEN
    event_name := 'job.' || NEW.status;
  ELSE
    RETURN NEW;
  END IF;

  PERFORM pg_advisory_xact_lock(hashtext(tenant));
  event_payload := jsonb_build_object(
    'job_id', NEW.id::text,
    'job_type', NEW.job_type,
    'status', NEW.status,
    'attempts', NEW.attempts,
    'max_attempts', NEW.max_attempts,
    'worker_id', COALESCE(NEW.locked_by, ''),
    'approval_id', COALESCE(NEW.approval_id::text, ''),
    'execution_action', COALESCE(NEW.execution_action, ''),
    'last_error', COALESCE(NEW.last_error, '')
  );

  INSERT INTO event_journal(
    tenant_id, event_type, sequence, correlation_id,
    causation_id, aggregate_id, payload
  ) VALUES (
    NEW.tenant_id, event_name,
    COALESCE((SELECT MAX(sequence)+1 FROM event_journal WHERE tenant_id=NEW.tenant_id), 1),
    corr, NULL, NEW.id::text, event_payload
  );

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_event_trigger ON durable_jobs;
CREATE TRIGGER durable_jobs_event_trigger
AFTER INSERT OR UPDATE OF status ON durable_jobs
FOR EACH ROW EXECUTE FUNCTION ftn_job_event_trigger();

CREATE INDEX IF NOT EXISTS event_journal_tenant_sequence_idx
  ON event_journal(tenant_id, sequence);
