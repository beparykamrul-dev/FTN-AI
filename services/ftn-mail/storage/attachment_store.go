package storage

import (
    "context"
    "crypto/sha256"
    "database/sql"
    "errors"
    "io"
)

type AttachmentMeta struct {
    ID string
    MessageID string
    MailboxID string
    Filename string
    MIMEType string
    SizeBytes int64
    SHA256 []byte
    StorageKey string
}

func (s *Store) GetAttachment(ctx context.Context, identityID, messageID, attachmentID string) (AttachmentMeta, error) {
    if s == nil || s.DB == nil || identityID == "" || messageID == "" || attachmentID == "" {
        return AttachmentMeta{}, errors.New("invalid attachment request")
    }
    var a AttachmentMeta
    err := s.DB.QueryRowContext(ctx, `SELECT a.id,a.message_id,a.mailbox_id,a.filename,a.mime_type,a.size_bytes,a.sha256,a.storage_key FROM ftn_mail.attachments a JOIN ftn_mailboxes b ON b.id=a.mailbox_id WHERE a.id=$1 AND a.message_id=$2 AND b.identity_id=$3`, attachmentID, messageID, identityID).Scan(&a.ID,&a.MessageID,&a.MailboxID,&a.Filename,&a.MIMEType,&a.SizeBytes,&a.SHA256,&a.StorageKey)
    if errors.Is(err, sql.ErrNoRows) { return AttachmentMeta{}, errors.New("attachment not found") }
    return a, err
}

func (s *Store) StreamAttachment(ctx context.Context, identityID, messageID, attachmentID string, blob *BlobStore, dst io.Writer) (AttachmentMeta, error) {
    if s == nil || s.DB == nil || blob == nil || dst == nil {
        return AttachmentMeta{}, errors.New("attachment stream unavailable")
    }
    a, err := s.GetAttachment(ctx, identityID, messageID, attachmentID)
    if err != nil { return AttachmentMeta{}, err }
    data, err := blob.Get(ctx, a.StorageKey)
    if err != nil { return AttachmentMeta{}, err }
    if int64(len(data)) != a.SizeBytes { return AttachmentMeta{}, errors.New("attachment size integrity check failed") }
    if len(a.SHA256) != sha256.Size { return AttachmentMeta{}, errors.New("attachment hash metadata invalid") }
    h := sha256.Sum256(data)
    if string(h[:]) != string(a.SHA256) { return AttachmentMeta{}, errors.New("attachment hash integrity check failed") }
    if _, err = dst.Write(data); err != nil { return AttachmentMeta{}, err }
    return a, nil
}
