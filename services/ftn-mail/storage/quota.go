package storage

import (
	"context"
	"database/sql"
	"errors"
)

var ErrQuotaExceeded = errors.New("mailbox quota exceeded")

func (s *Store) CheckQuota(ctx context.Context, mailboxID string, incomingBytes int64) error {
	if mailboxID == "" || incomingBytes < 0 {
		return errors.New("invalid quota request")
	}
	var quota, used int64
	err := s.DB.QueryRowContext(ctx, `
		SELECT m.quota_bytes,
		       COALESCE((SELECT SUM(size_bytes) FROM ftn_mail_messages WHERE mailbox_id=m.id),0)
		FROM ftn_mailboxes m
		WHERE m.id=$1 AND m.status='active'`, mailboxID).Scan(&quota, &used)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { return errors.New("mailbox unavailable") }
		return err
	}
	if incomingBytes > quota-used { return ErrQuotaExceeded }
	return nil
}
