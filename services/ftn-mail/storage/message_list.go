package storage

import (
	"context"
	"errors"
)

type MessageListItem struct {
	ID string `json:"id"`
	Sender string `json:"sender"`
	Subject string `json:"subject"`
	SizeBytes int64 `json:"size_bytes"`
	CreatedAt string `json:"created_at"`
}

func (s *Store) ListMessages(ctx context.Context, identityID, localPart, domain string, limit, offset int) ([]MessageListItem, error) {
	if s == nil || s.DB == nil || identityID == "" || localPart == "" || domain == "" { return nil, errors.New("invalid mailbox query") }
	if limit <= 0 || limit > 100 { limit = 50 }
	if offset < 0 { offset = 0 }
	rows, err := s.DB.QueryContext(ctx, `SELECT m.id, COALESCE(m.sender,''), COALESCE(m.subject,''), m.size_bytes, m.created_at FROM ftn_mail_messages m JOIN ftn_mailboxes b ON b.id=m.mailbox_id WHERE b.identity_id=$1 AND b.local_part=$2 AND b.domain=$3 ORDER BY m.created_at DESC LIMIT $4 OFFSET $5`, identityID, localPart, domain, limit, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	out := make([]MessageListItem, 0)
	for rows.Next() {
		var m MessageListItem
		if err := rows.Scan(&m.ID,&m.Sender,&m.Subject,&m.SizeBytes,&m.CreatedAt); err != nil { return nil, err }
		out = append(out,m)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return out,nil
}
