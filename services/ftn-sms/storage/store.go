package storage

import (
	"context"
	"database/sql"
	"errors"
)

type Message struct {
	ID          string
	IdentityID  string
	SenderID    string
	Recipient   string
	Body        string
	Priority    int
	ScheduledAt sql.NullTime
	Attempts    int
}

type Store struct{ DB *sql.DB }

func (s *Store) Enqueue(ctx context.Context, m Message) error {
	if s == nil || s.DB == nil || m.ID == "" || m.IdentityID == "" || m.SenderID == "" || m.Recipient == "" || m.Body == "" {
		return errors.New("invalid sms message")
	}
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO ftn_sms_messages
		(id, identity_id, sender_id, recipient, body, priority, scheduled_at, status, attempts)
		VALUES ($1,$2,$3,$4,$5,$6,COALESCE($7,NOW()),'QUEUED',0)`,
		m.ID, m.IdentityID, m.SenderID, m.Recipient, m.Body, m.Priority, m.ScheduledAt)
	return err
}

func (s *Store) Claim(ctx context.Context) (Message, error) {
	if s == nil || s.DB == nil {
		return Message{}, errors.New("sms store unavailable")
	}
	var m Message
	err := s.DB.QueryRowContext(ctx, `
		WITH next_message AS (
			SELECT id FROM ftn_sms_messages
			WHERE status='QUEUED' AND scheduled_at <= NOW()
			ORDER BY priority DESC, scheduled_at ASC, created_at ASC
			FOR UPDATE SKIP LOCKED LIMIT 1
		)
		UPDATE ftn_sms_messages m
		SET status='PROCESSING', updated_at=NOW()
		FROM next_message n
		WHERE m.id=n.id
		RETURNING m.id,m.identity_id,m.sender_id,m.recipient,m.body,m.priority,m.scheduled_at,m.attempts`,
	).Scan(&m.ID, &m.IdentityID, &m.SenderID, &m.Recipient, &m.Body, &m.Priority, &m.ScheduledAt, &m.Attempts)
	return m, err
}

func (s *Store) MarkSent(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE ftn_sms_messages SET status='SENT', updated_at=NOW() WHERE id=$1 AND status='PROCESSING'`, id)
	return err
}

func (s *Store) MarkDelivered(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE ftn_sms_messages SET status='DELIVERED', updated_at=NOW() WHERE id=$1 AND status IN ('SENT','PROCESSING')`, id)
	return err
}

func (s *Store) Requeue(ctx context.Context, id string, retryAt string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE ftn_sms_messages SET status='QUEUED', attempts=attempts+1, scheduled_at=$2, updated_at=NOW() WHERE id=$1 AND status='PROCESSING'`, id, retryAt)
	return err
}

func (s *Store) MarkFailed(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE ftn_sms_messages SET status='FAILED', attempts=attempts+1, updated_at=NOW() WHERE id=$1 AND status='PROCESSING'`, id)
	return err
}
