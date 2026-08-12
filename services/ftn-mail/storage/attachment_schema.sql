CREATE TABLE IF NOT EXISTS ftn_mail.attachments (
    id UUID PRIMARY KEY,
    message_id UUID NOT NULL,
    mailbox_id UUID NOT NULL,
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    sha256 BYTEA NOT NULL CHECK (octet_length(sha256) = 32),
    storage_key TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_attachments_message
ON ftn_mail.attachments(message_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_attachments_mailbox
ON ftn_mail.attachments(mailbox_id, created_at DESC);
