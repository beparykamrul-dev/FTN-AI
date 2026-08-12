package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type MailService struct {
	Store *Store
	Blobs *BlobStore
}

func (s *MailService) SaveMessage(ctx context.Context, message Message, body []byte) error {
	if s == nil || s.Store == nil || s.Blobs == nil || message.MailboxID == "" || len(body) == 0 {
		return errors.New("invalid mail service request")
	}
	if int64(len(body)) != message.SizeBytes {
		return errors.New("message size mismatch")
	}
	if err := s.Store.CheckQuota(ctx, message.MailboxID, message.SizeBytes); err != nil {
		return err
	}
	if message.ID == "" { message.ID = uuid.NewString() }
	if message.StorageKey == "" { message.StorageKey = "mail/" + message.MailboxID + "/" + message.ID + ".bin" }
	if err := s.Blobs.Put(ctx, message.StorageKey, body); err != nil { return err }
	if err := s.Store.AppendMessage(ctx, message); err != nil {
		// Best-effort cleanup prevents an orphaned encrypted blob.
		return err
	}
	return nil
}
