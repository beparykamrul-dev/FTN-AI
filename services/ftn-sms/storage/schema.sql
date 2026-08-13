-- FTN SMS service V1 persistence contract.
-- No third-party provider state is stored here.
-- Actual sender IDs must be operator-approved.

CREATE TABLE IF NOT EXISTS sms_sender_ids (
    id UUID PRIMARY KEY,
    sender VARCHAR(32) NOT NULL UNIQUE,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (status IN ('pending', 'approved', 'suspended', 'retired'))
);

CREATE TABLE IF NOT EXISTS sms_messages (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL,
    sender_id UUID NOT NULL REFERENCES sms_sender_ids(id),
    recipient VARCHAR(32) NOT NULL,
    body TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'queued',
    priority SMALLINT NOT NULL DEFAULT 0,
    attempts INTEGER NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (status IN ('queued', 'processing', 'sent', 'delivered', 'failed', 'cancelled')),
    CHECK (priority BETWEEN 0 AND 9),
    CHECK (attempts >= 0)
);

CREATE INDEX IF NOT EXISTS idx_sms_messages_identity_created
    ON sms_messages(identity_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_sms_messages_queue
    ON sms_messages(status, priority DESC, scheduled_at, created_at);

CREATE TABLE IF NOT EXISTS sms_delivery_events (
    id BIGSERIAL PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES sms_messages(id) ON DELETE CASCADE,
    event_type VARCHAR(32) NOT NULL,
    operator_status VARCHAR(64),
    provider_message_id VARCHAR(128),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sms_delivery_events_message
    ON sms_delivery_events(message_id, occurred_at DESC);
