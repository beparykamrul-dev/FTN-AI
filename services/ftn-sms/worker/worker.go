package worker

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/beparykamrul-dev/FTN-AI/services/ftn-sms/storage"
	"github.com/beparykamrul-dev/FTN-AI/services/ftn-sms/transport"
)

type Worker struct {
	Store     *storage.Store
	Transport transport.Adapter
	MaxRetry  int
}

func (w *Worker) RunOnce(ctx context.Context) error {
	if w == nil || w.Store == nil || w.Transport == nil {
		return errors.New("sms worker not configured")
	}
	if !w.Transport.Ready(ctx) {
		return errors.New("sms transport not ready")
	}

	m, err := w.Store.Claim(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	result, sendErr := w.Transport.Send(ctx, transport.Message{
		ID: m.ID, Sender: m.SenderID, Recipient: m.Recipient, Body: m.Body,
	})
	if sendErr == nil {
		switch result {
		case transport.Sent:
			return w.Store.MarkSent(ctx, m.ID)
		case transport.Delivered:
			return w.Store.MarkDelivered(ctx, m.ID)
		case transport.PermanentFailure:
			return w.Store.MarkFailed(ctx, m.ID)
		case transport.TemporaryFailure:
			return w.retryOrFail(ctx, m)
		default:
			return w.Store.MarkFailed(ctx, m.ID)
		}
	}

	if result == transport.TemporaryFailure && m.Attempts < w.MaxRetry {
		return w.retryOrFail(ctx, m)
	}
	return w.Store.MarkFailed(ctx, m.ID)
}

func (w *Worker) retryOrFail(ctx context.Context, m storage.Message) error {
	if w.MaxRetry > 0 && m.Attempts >= w.MaxRetry {
		return w.Store.MarkFailed(ctx, m.ID)
	}
	backoff := time.Second * time.Duration(1<<min(m.Attempts, 6))
	return w.Store.Requeue(ctx, m.ID, time.Now().Add(backoff).UTC().Format(time.RFC3339Nano))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
