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
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, identity_id, local_part, domain, status, quota_bytes
		FROM ftn_mailboxes
		WHERE identity_id=$1 AND local_part=$2 AND domain=$3`, identityID, localPart, domain).
		Scan(&m.ID, &m.IdentityID, &m.LocalPart, &m.Domain, &m.Status, &m.QuotaBytes)
	return m, err
}

func (s *Store) AppendMessage(ctx context.Context, m Message) error {
	if m.MailboxID == "" || m.StorageKey == "" || m.SizeBytes < 0 {
		return errors.New("invalid message")
	}
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO ftn_mail_messages
		(id, mailbox_id, message_uid, message_id, sender, recipients, subject, size_bytes, storage_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		m.ID, m.MailboxID, m.UID, m.MessageID, m.Sender, m.Recipients, m.Subject, m.SizeBytes, m.StorageKey)
	return err
}
