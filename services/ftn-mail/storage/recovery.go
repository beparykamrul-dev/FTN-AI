package storage

import (
	"context"
	"database/sql"
)

type OutboxRecord struct {
	ID         string
	MessageID  string
	StorageKey string
	State      string
	Attempts   int
}

func PendingBlobOutbox(ctx context.Context, db *sql.DB, limit int) ([]OutboxRecord, error) {
	if limit <= 0 || limit > 500 { limit = 100 }
	rows, err := db.QueryContext(ctx, `
		SELECT id, message_id, storage_key, state, attempts
		FROM ftn_mail_blob_outbox
		WHERE state IN ('pending','cleanup')
		ORDER BY updated_at ASC
		LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []OutboxRecord
	for rows.Next() {
		var r OutboxRecord
		if err := rows.Scan(&r.ID, &r.MessageID, &r.StorageKey, &r.State, &r.Attempts); err != nil { return nil, err }
		result = append(result, r)
	}
	return result, rows.Err()
}
