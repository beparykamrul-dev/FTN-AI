package storage

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Store) MarkMessageRead(ctx context.Context, messageID, identityID string) (bool, error) {
	if s == nil || s.DB == nil || messageID == "" || identityID == "" { return false, errors.New("invalid read-state request") }
	var changed bool
	err := s.DB.QueryRowContext(ctx, `SELECT ftn_mail_mark_message_read($1,$2)`, messageID, identityID).Scan(&changed)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { return false, err }
	return changed, nil
}

func (s *Store) UnreadCount(ctx context.Context, identityID, localPart, domain string) (int64, error) {
	if s == nil || s.DB == nil || identityID == "" || localPart == "" || domain == "" { return 0, errors.New("invalid unread-count request") }
	var count int64
	if err := s.DB.QueryRowContext(ctx, `SELECT ftn_mail_unread_count($1,$2,$3)`, identityID, localPart, domain).Scan(&count); err != nil { return 0, err }
	return count, nil
}
