-- Restore the approval lifecycle event bridge from the legacy duplicate 006 migration.
-- The canonical job lifecycle trigger lives in 006_job_event_automation.sql;
-- this migration installs only the approval-side trigger.
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

INSERT INTO schema_migrations(version,name)
VALUES (15,'015_approval_event_trigger.sql')
ON CONFLICT (version) DO NOTHING;
