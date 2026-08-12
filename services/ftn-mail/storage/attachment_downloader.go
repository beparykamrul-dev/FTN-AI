package storage

import (
	"context"
	"errors"
	"io"
)

// AttachmentBlobReader is the narrow storage contract required by HTTP/Android/Web adapters.
type AttachmentBlobReader interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

type AttachmentDownloader struct {
	Store *Store
	Blob  AttachmentBlobReader
}

func (d *AttachmentDownloader) Stream(ctx context.Context, identityID, messageID, attachmentID string, dst io.Writer) (AttachmentMeta, error) {
	if d == nil || d.Store == nil || d.Blob == nil {
		return AttachmentMeta{}, errors.New("attachment downloader unavailable")
	}
	a, err := d.Store.GetAttachment(ctx, identityID, messageID, attachmentID)
	if err != nil {
		return AttachmentMeta{}, err
	}
	data, err := d.Blob.Get(ctx, a.StorageKey)
	if err != nil {
		return AttachmentMeta{}, err
	}
	if int64(len(data)) != a.SizeBytes {
		return AttachmentMeta{}, errors.New("attachment size integrity check failed")
	}
	if _, err := dst.Write(data); err != nil {
		return AttachmentMeta{}, err
	}
	return a, nil
}
