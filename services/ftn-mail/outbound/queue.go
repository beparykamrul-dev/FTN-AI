package outbound

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Queue struct { DB *sql.DB }

func (q *Queue) Enqueue(ctx context.Context, sender string, recipients []string, raw []byte) (string, error) {
	if q == nil || q.DB == nil || sender == "" || len(recipients) == 0 || len(raw) == 0 { return "", errors.New("invalid outbound message") }
	id := uuid.NewString()
	_, err := q.DB.ExecContext(ctx, `INSERT INTO ftn_mail_outbound_queue (id,sender,recipients,raw_message,status,next_attempt_at,created_at,updated_at) VALUES ($1,$2,$3,$4,'queued',$5,now(),now())`, id, sender, recipients, raw, time.Now().UTC())
	if err != nil { return "", err }
	return id, nil
}
