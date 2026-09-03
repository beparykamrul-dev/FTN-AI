-- FTN approval-gated execution boundary. Additive/idempotent.
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS approval_id UUID REFERENCES change_approvals(id) ON DELETE SET NULL;
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS approval_request_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS execution_action TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS rollback_payload JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS verification_payload JSONB NOT NULL DEFAULT '{}'::jsonb;
CREATE INDEX IF NOT EXISTS durable_jobs_approval_idx ON durable_jobs(approval_id);
CREATE INDEX IF NOT EXISTS durable_jobs_approval_request_hash_idx ON durable_jobs(approval_request_hash);
CREATE INDEX IF NOT EXISTS change_approvals_request_hash_idx ON change_approvals(request_hash);

CREATE OR REPLACE FUNCTION ftn_bind_durable_job_approval()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    approval_tenant UUID;
    approval_action TEXT;
    approval_hash TEXT;
BEGIN
    IF NEW.approval_id IS NULL THEN
        NEW.approval_request_hash := '';
        RETURN NEW;
    END IF;

    SELECT tenant_id, action, request_hash
      INTO approval_tenant, approval_action, approval_hash
      FROM change_approvals
     WHERE id = NEW.approval_id
     FOR SHARE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'approval_not_found';
    END IF;
    IF approval_tenant IS DISTINCT FROM NEW.tenant_id THEN
        RAISE EXCEPTION 'approval_tenant_mismatch';
    END IF;
    IF NULLIF(BTRIM(NEW.execution_action), '') IS NULL THEN
        RAISE EXCEPTION 'execution_action_required_for_approved_job';
    END IF;
    IF approval_action IS DISTINCT FROM NEW.execution_action THEN
        RAISE EXCEPTION 'approval_action_mismatch';
    END IF;

    NEW.approval_request_hash := approval_hash;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS durable_jobs_approval_binding ON durable_jobs;
CREATE TRIGGER durable_jobs_approval_binding
BEFORE INSERT OR UPDATE OF approval_id, execution_action, tenant_id
ON durable_jobs
FOR EACH ROW
EXECUTE FUNCTION ftn_bind_durable_job_approval();

INSERT INTO permissions(key, description) VALUES
 ('approval.create','Create privileged change approvals'),
 ('approval.decide','Approve or reject privileged changes'),
 ('job.verify','Verify execution results'),
 ('job.rollback','Rollback an executed change')
ON CONFLICT (key) DO NOTHING;
