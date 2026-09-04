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

	id := uuid.NewString()
	h := sha256.Sum256(raw)
	storageKey := fmt.Sprintf("mail/%s/%s-%x.bin", mailboxID, id, h[:8])
	if err := s.Blobs.Put(ctx, storageKey, raw); err != nil { return err }

	message := Message{
		ID: id,
		MailboxID: mailboxID,
		MessageID: "<" + id + "@" + domain + ">",
		Sender: sender,
		Recipients: recipients,
		SizeBytes: int64(len(raw)),
		StorageKey: storageKey,
	}
	return (&Store{DB: s.DB}).AppendMessageWithQuota(ctx, message)
}
