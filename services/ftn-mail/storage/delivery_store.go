package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type DeliveryStore struct {
	DB    *sql.DB
	Blobs *BlobStore
}

func (s *DeliveryStore) Deliver(ctx context.Context, identityID, localPart, domain, sender string, recipients []string, raw []byte) error {
	if s == nil || s.DB == nil || s.Blobs == nil || identityID == "" || len(raw) == 0 {
		return fmt.Errorf("invalid delivery store")
	}
	var mailboxID string
	err := s.DB.QueryRowContext(ctx, `
		SELECT id FROM ftn_mailboxes
		WHERE identity_id=$1 AND local_part=$2 AND domain=$3 AND status='active'`,
		identityID, strings.ToLower(localPart), strings.ToLower(domain)).Scan(&mailboxID)
	if err != nil { return err }

	if err := (&Store{DB: s.DB}).CheckQuota(ctx, mailboxID, int64(len(raw))); err != nil { return err }

	id := uuid.NewString()
	h := sha256.Sum256(raw)
	storageKey := fmt.Sprintf("mail/%s/%s-%x.bin", mailboxID, id, h[:8])
	if err := s.Blobs.Put(ctx, storageKey, raw); err != nil { return err }

	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO ftn_mail_messages
		(id, mailbox_id, message_uid, message_id, sender, recipients, size_bytes, storage_key)
		VALUES ($1,$2,COALESCE((SELECT MAX(message_uid)+1 FROM ftn_mail_messages WHERE mailbox_id=$2),1),$3,$4,$5,$6,$7)`,
		id, mailboxID, "<"+id+"@"+domain+">", sender, recipients, len(raw), storageKey)
	if err != nil { return err }
	return nil
}
