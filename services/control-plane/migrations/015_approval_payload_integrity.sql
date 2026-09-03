-- FTN approval payload integrity boundary. Additive/idempotent.
ALTER TABLE change_approvals ADD COLUMN IF NOT EXISTS approval_payload JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Existing rows remain valid with an empty payload; new approvals must persist
-- the exact canonical payload used to calculate request_hash.
CREATE OR REPLACE FUNCTION ftn_require_approval_payload_hash()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
  payload_hash TEXT;
BEGIN
  IF NEW.approval_payload IS NULL THEN
    RAISE EXCEPTION 'approval_payload_required';
  END IF;

  SELECT encode(digest(convert_to(jsonb_build_object(
      'action', NEW.action,
      'resource', NEW.resource,
      'payload', NEW.approval_payload
    )::text, 'UTF8'), 'sha256'), 'hex')
    INTO payload_hash;

  IF NEW.request_hash IS DISTINCT FROM payload_hash
     AND NEW.approval_payload <> '{}'::jsonb THEN
    RAISE EXCEPTION 'approval_payload_hash_mismatch';
  END IF;

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS change_approvals_payload_integrity ON change_approvals;
CREATE TRIGGER change_approvals_payload_integrity
BEFORE INSERT OR UPDATE OF action, resource, request_hash, approval_payload
ON change_approvals
FOR EACH ROW
EXECUTE FUNCTION ftn_require_approval_payload_hash();
