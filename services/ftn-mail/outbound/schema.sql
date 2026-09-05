CREATE TABLE IF NOT EXISTS ftn_mail_outbound_queue (
    id UUID PRIMARY KEY,
    sender TEXT NOT NULL,
    recipients TEXT[] NOT NULL,
    raw_message BYTEA NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','processing','retry','sent','failed')),
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_error TEXT,
    lease_token UUID,
    lease_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_outbound_ready
    ON ftn_mail_outbound_queue(status, next_attempt_at);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_outbound_processing_lease
    ON ftn_mail_outbound_queue(status, lease_expires_at)
    WHERE status = 'processing';
