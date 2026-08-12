-- Mail read-state migration for shared FTN PostgreSQL.
ALTER TABLE ftn_mail_messages
    ADD COLUMN IF NOT EXISTS is_read BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_ftn_mail_messages_unread
    ON ftn_mail_messages(mailbox_id, created_at DESC)
    WHERE is_read = FALSE;

CREATE OR REPLACE FUNCTION ftn_mail_mark_message_read(p_message_id UUID, p_identity_id UUID)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE changed BOOLEAN;
BEGIN
    UPDATE ftn_mail_messages m
       SET is_read = TRUE
      FROM ftn_mailboxes b
     WHERE m.id = p_message_id
       AND m.mailbox_id = b.id
       AND b.identity_id = p_identity_id
       AND m.is_read = FALSE;
    changed := FOUND;
    RETURN changed;
END;
$$;

CREATE OR REPLACE FUNCTION ftn_mail_unread_count(p_identity_id UUID, p_local_part TEXT, p_domain TEXT)
RETURNS BIGINT
LANGUAGE SQL STABLE
AS $$
    SELECT COUNT(*)
      FROM ftn_mail_messages m
      JOIN ftn_mailboxes b ON b.id = m.mailbox_id
     WHERE b.identity_id = p_identity_id
       AND b.local_part = p_local_part
       AND b.domain = p_domain
       AND m.is_read = FALSE;
$$;
