CREATE SCHEMA IF NOT EXISTS ftn_mail;

CREATE TABLE IF NOT EXISTS ftn_mail.outbound_queue (
    id UUID PRIMARY KEY,
    sender TEXT NOT NULL,
    recipients TEXT[] NOT NULL,
    raw_message BYTEA NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','processing','retry','sent','failed')),
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_outbound_ready
ON ftn_mail.outbound_queue(status, next_attempt_at);

CREATE TABLE IF NOT EXISTS ftn_mail.delivery_events (
    id BIGSERIAL PRIMARY KEY,
    queue_id UUID REFERENCES ftn_mail.outbound_queue(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    status_code TEXT,
    action TEXT,
    diagnostic TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_delivery_events_queue
ON ftn_mail.delivery_events(queue_id, created_at DESC);
