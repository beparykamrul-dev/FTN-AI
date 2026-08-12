CREATE TABLE IF NOT EXISTS ftn_mail_blob_outbox (
    id UUID PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES ftn_mail_messages(id) ON DELETE CASCADE,
    storage_key TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'pending' CHECK (state IN ('pending','committed','cleanup')),
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_blob_outbox_pending
    ON ftn_mail_blob_outbox(state, updated_at);
