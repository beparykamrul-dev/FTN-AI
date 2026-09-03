-- FTN approval payload integrity boundary. Additive/idempotent.
ALTER TABLE change_approvals ADD COLUMN IF NOT EXISTS approval_payload JSONB NOT NULL DEFAULT '{}'::jsonb;

-- The application computes request_hash from the exact canonical approval request.
-- Persisting the payload makes that reviewed object durable and prevents it from
-- being silently replaced after approval.
CREATE OR REPLACE FUNCTION ftn_require_approval_payload_integrity()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.approval_payload IS NULL THEN
    RAISE EXCEPTION 'approval_payload_required';
  END IF;

  IF TG_OP = 'UPDATE' AND OLD.status IN ('approved','executed','rolled_back') THEN
    IF NEW.action IS DISTINCT FROM OLD.action
       OR NEW.resource IS DISTINCT FROM OLD.resource
       OR NEW.request_hash IS DISTINCT FROM OLD.request_hash
       OR NEW.approval_payload IS DISTINCT FROM OLD.approval_payload THEN
      RAISE EXCEPTION 'approved_request_immutable';
    END IF;
  END IF;

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_payload_integrity ON change_approvals;
CREATE TRIGGER change_approvals_payload_integrity
BEFORE INSERT OR UPDATE OF action, resource, request_hash, approval_payload
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_require_approval_payload_integrity();
