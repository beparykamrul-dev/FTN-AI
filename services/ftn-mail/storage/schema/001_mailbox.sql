CREATE TABLE IF NOT EXISTS ftn_mailboxes (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL,
    local_part TEXT NOT NULL,
    domain TEXT NOT NULL DEFAULT 'familytimenet.com',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled','suspended')),
    quota_bytes BIGINT NOT NULL DEFAULT 1073741824 CHECK (quota_bytes > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(local_part, domain)
);

CREATE INDEX IF NOT EXISTS idx_ftn_mailboxes_identity
    ON ftn_mailboxes(identity_id, status);

CREATE TABLE IF NOT EXISTS ftn_mail_messages (
    id UUID PRIMARY KEY,
    mailbox_id UUID NOT NULL REFERENCES ftn_mailboxes(id) ON DELETE CASCADE,
    message_uid BIGINT NOT NULL,
    message_id TEXT,
    sender TEXT NOT NULL,
    recipients TEXT[] NOT NULL,
    subject TEXT,
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    storage_key TEXT NOT NULL UNIQUE,
    flags TEXT[] NOT NULL DEFAULT '{}',
    UNIQUE(mailbox_id, message_uid)
);

CREATE INDEX IF NOT EXISTS idx_ftn_mail_messages_mailbox_date
    ON ftn_mail_messages(mailbox_id, received_at DESC);
