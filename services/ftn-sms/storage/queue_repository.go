package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Message is the persisted SMS work item claimed by a worker.
type Message struct {
	ID          string
	IdentityID  string
	SenderID    string
	Recipient   string
	Body        string
	Priority    int
	Attempts    int
	ScheduledAt *time.Time
}

type QueueRepository struct {
	DB *sql.DB
}

func (r *QueueRepository) Enqueue(ctx context.Context, identityID, senderID, recipient, body string, priority int, scheduledAt *time.Time) (string, error) {
	if r == nil || r.DB == nil {
		return "", errors.New("sms: database is nil")
	}
	if identityID == "" || senderID == "" || recipient == "" || body == "" {
		return "", errors.New("sms: required field is empty")
	}
	if priority < 0 || priority > 9 {
		return "", errors.New("sms: priority out of range")
	}

	id, err := newUUID()
	if err != nil {
		return "", fmt.Errorf("sms: generate message id: %w", err)
	}

	_, err = r.DB.ExecContext(ctx, `
		INSERT INTO sms_messages
			(id, identity_id, sender_id, recipient, body, status, priority, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'queued', $6, $7, NOW(), NOW())
	`, id, identityID, senderID, recipient, body, priority, scheduledAt)
	if err != nil {
		return "", fmt.Errorf("sms: enqueue: %w", err)
	}
	return id, nil
}

// Claim atomically leases one eligible queued message to this worker.
// The lease is represented by the processing state; recovery of stale processing
// records is handled by the recovery operation below.
func (r *QueueRepository) Claim(ctx context.Context) (*Message, error) {
	if r == nil || r.DB == nil {
		return nil, errors.New("sms: database is nil")
	}

	row := r.DB.QueryRowContext(ctx, `
		WITH candidate AS (
			SELECT id
			FROM sms_messages
			WHERE status = 'queued'
			  AND (scheduled_at IS NULL OR scheduled_at <= NOW())
			ORDER BY priority DESC, scheduled_at NULLS FIRST, created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		UPDATE sms_messages m
		SET status = 'processing', updated_at = NOW()
		FROM candidate c
		WHERE m.id = c.id
		RETURNING m.id, m.identity_id, m.sender_id, m.recipient, m.body,
		          m.priority, m.attempts, m.scheduled_at
	`)

	var m Message
	if err := row.Scan(&m.ID, &m.IdentityID, &m.SenderID, &m.Recipient, &m.Body, &m.Priority, &m.Attempts, &m.ScheduledAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("sms: claim: %w", err)
	}
	return &m, nil
}

func (r *QueueRepository) MarkSent(ctx context.Context, id string) error {
	return r.transition(ctx, id, "processing", "sent")
}

func (r *QueueRepository) MarkDelivered(ctx context.Context, id string) error {
	return r.transition(ctx, id, "sent", "delivered")
}

func (r *QueueRepository) MarkFailed(ctx context.Context, id string) error {
	return r.transition(ctx, id, "processing", "failed")
}

func (r *QueueRepository) Requeue(ctx context.Context, id string, nextAttempt time.Time) error {
	if r == nil || r.DB == nil {
		return errors.New("sms: database is nil")
	}
	result, err := r.DB.ExecContext(ctx, `
		UPDATE sms_messages
		SET status = 'queued', attempts = attempts + 1, scheduled_at = $2, updated_at = NOW()
		WHERE id = $1 AND status = 'processing'
	`, id, nextAttempt)
	if err != nil {
		return fmt.Errorf("sms: requeue: %w", err)
	}
	return requireOne(result, "sms: requeue")
}

func (r *QueueRepository) transition(ctx context.Context, id, from, to string) error {
	if r == nil || r.DB == nil {
		return errors.New("sms: database is nil")
	}
	result, err := r.DB.ExecContext(ctx, `
		UPDATE sms_messages SET status = $2, updated_at = NOW()
		WHERE id = $1 AND status = $3
	`, id, to, from)
	if err != nil {
		return fmt.Errorf("sms: transition %s->%s: %w", from, to, err)
	}
	return requireOne(result, "sms: state transition")
}

func requireOne(result sql.Result, operation string) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", operation, err)
	}
	if rows != 1 {
		return fmt.Errorf("%s: message not in expected state", operation)
	}
	return nil
}

func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], b[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], b[10:16])
	return string(buf), nil
}
