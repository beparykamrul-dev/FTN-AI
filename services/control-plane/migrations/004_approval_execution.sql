-- FTN approval-gated execution boundary. Additive/idempotent.
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS approval_id UUID REFERENCES change_approvals(id) ON DELETE SET NULL;
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS execution_action TEXT NOT NULL DEFAULT '';
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS rollback_payload JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE durable_jobs ADD COLUMN IF NOT EXISTS verification_payload JSONB NOT NULL DEFAULT '{}'::jsonb;
CREATE INDEX IF NOT EXISTS durable_jobs_approval_idx ON durable_jobs(approval_id);
CREATE INDEX IF NOT EXISTS change_approvals_request_hash_idx ON change_approvals(request_hash);
INSERT INTO permissions(key, description) VALUES
 ('approval.create','Create privileged change approvals'),
 ('approval.decide','Approve or reject privileged changes'),
 ('job.verify','Verify execution results'),
 ('job.rollback','Rollback an executed change')
ON CONFLICT (key) DO NOTHING;
