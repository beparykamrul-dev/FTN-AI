-- FTN verification integrity boundary. Additive/idempotent.
-- An approval can transition to executed only when its exact bound durable job
-- has completed successfully and carries a successful verification envelope.
CREATE OR REPLACE FUNCTION ftn_require_verified_approval_execution()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
  verified_count INTEGER;
BEGIN
  IF NEW.status IS DISTINCT FROM 'executed' OR OLD.status IS NOT DISTINCT FROM 'executed' THEN
    RETURN NEW;
  END IF;

  SELECT COUNT(*)
    INTO verified_count
    FROM durable_jobs j
   WHERE j.approval_id = NEW.id
     AND j.tenant_id IS NOT DISTINCT FROM NEW.tenant_id
     AND j.status = 'succeeded'
     AND j.execution_action = NEW.action
     AND j.approval_request_hash = NEW.request_hash
     AND j.verification_payload->>'success' = 'true';

  IF verified_count <> 1 THEN
    RAISE EXCEPTION 'approval_execution_requires_exact_verified_job';
  END IF;

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_verified_execution ON change_approvals;
CREATE TRIGGER change_approvals_verified_execution
BEFORE UPDATE OF status
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_require_verified_approval_execution();
