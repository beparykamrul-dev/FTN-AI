package storage

import (
	"context"
	"database/sql"
	"errors"
)

type Mailbox struct {
	ID         string
	IdentityID string
	LocalPart  string
	Domain     string
	Status     string
	QuotaBytes int64
}

type Message struct {
	ID         string
	MailboxID  string
	UID        int64
	MessageID  string
	Sender     string
	Recipients []string
	Subject    string
	SizeBytes  int64
	StorageKey string
}

type Store struct{ DB *sql.DB }

func (s *Store) GetMailbox(ctx context.Context, identityID, localPart, domain string) (Mailbox, error) {
	var m Mailbox
	if s == nil || s.DB == nil || identityID == "" || localPart == "" || domain == "" {
		return m, errors.New("invalid mailbox lookup")
	}
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, identity_id, local_part, domain, status, quota_bytes
		FROM ftn_mailboxes
		WHERE identity_id=$1 AND local_part=$2 AND domain=$3`, identityID, localPart, domain).
		Scan(&m.ID, &m.IdentityID, &m.LocalPart, &m.Domain, &m.Status, &m.QuotaBytes)
	return m, err
}

func (s *Store) AppendMessage(ctx context.Context, m Message) error {
	if s == nil || s.DB == nil || m.MailboxID == "" || m.StorageKey == "" || m.SizeBytes < 0 {
		return errors.New("invalid message")
	}
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO ftn_mail_messages
		(id, mailbox_id, message_uid, message_id, sender, recipients, subject, size_bytes, storage_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		m.ID, m.MailboxID, m.UID, m.MessageID, m.Sender, m.Recipients, m.Subject, m.SizeBytes, m.StorageKey)
	return err
}

// AppendMessageWithQuota serializes quota check and message insertion per mailbox.
// The advisory transaction lock closes the check-then-insert race without requiring
// a denormalized used-bytes counter.
func (s *Store) AppendMessageWithQuota(ctx context.Context, m Message) error {
	if s == nil || s.DB == nil || m.MailboxID == "" || m.StorageKey == "" || m.SizeBytes < 0 {
		return errors.New("invalid message")
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil { return err }
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1::text,0))`, m.MailboxID); err != nil {
		return err
	}
	var quota, used int64
	if err = tx.QueryRowContext(ctx, `
		SELECT m.quota_bytes,
		       COALESCE((SELECT SUM(size_bytes) FROM ftn_mail_messages WHERE mailbox_id=m.id),0)
		FROM ftn_mailboxes m
		WHERE m.id=$1 AND m.status='active'`, m.MailboxID).Scan(&quota, &used); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return errors.New("mailbox unavailable") }
		return err
	}
	if m.SizeBytes > quota-used { return ErrQuotaExceeded }

	if m.UID <= 0 {
		if err = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(message_uid)+1,1) FROM ftn_mail_messages WHERE mailbox_id=$1`, m.MailboxID).Scan(&m.UID); err != nil {
			return err
		}
	}
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO ftn_mail_messages
		(id, mailbox_id, message_uid, message_id, sender, recipients, subject, size_bytes, storage_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		m.ID, m.MailboxID, m.UID, m.MessageID, m.Sender, m.Recipients, m.Subject, m.SizeBytes, m.StorageKey); err != nil {
		return err
	}
	return tx.Commit()
}
