-- Make outbound delivery ownership explicit so stale recovery cannot
-- complete or overwrite a job claimed by another worker.
ALTER TABLE ftn_mail_outbound_queue
    ADD COLUMN IF NOT EXISTS lease_token UUID,
    ADD COLUMN IF NOT EXISTS lease_expires_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_ftn_mail_outbound_processing_lease
    ON ftn_mail_outbound_queue(status, lease_expires_at)
    WHERE status = 'processing';

INSERT INTO schema_migrations(version,name)
VALUES (24,'024_outbound_queue_leases.sql')
ON CONFLICT (version) DO NOTHING;
