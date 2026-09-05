package agent

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Plan defines the usage entitlement for public and customer AI access.
type Plan struct {
	ID              string
	RequestsPerDay  int
	TokensPerDay    int64
	MaxInputTokens  int64
	MaxOutputTokens int64
}

var Plans = map[string]Plan{
	"free":     {ID: "free", RequestsPerDay: 20, TokensPerDay: 20000, MaxInputTokens: 4000, MaxOutputTokens: 2000},
	"basic":    {ID: "basic", RequestsPerDay: 100, TokensPerDay: 100000, MaxInputTokens: 8000, MaxOutputTokens: 4000},
	"standard": {ID: "standard", RequestsPerDay: 500, TokensPerDay: 500000, MaxInputTokens: 16000, MaxOutputTokens: 8000},
	"pro":      {ID: "pro", RequestsPerDay: 2000, TokensPerDay: 2000000, MaxInputTokens: 32000, MaxOutputTokens: 16000},
}

type Usage struct {
	Requests int
	Tokens   int64
	ResetAt  time.Time
}

type QuotaStore interface {
	Get(ctx context.Context, scope Scope) (Usage, error)
	Put(ctx context.Context, scope Scope, usage Usage) error
}

type QuotaGate struct {
	Store QuotaStore
	mu    sync.Mutex
}

func (q *QuotaGate) CheckAndConsume(ctx context.Context, scope Scope, plan Plan, tokens int64) error {
	if q == nil || q.Store == nil {
		return errors.New("quota store is required")
	}
	if ctx == nil {
		return errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if plan.RequestsPerDay <= 0 || plan.TokensPerDay <= 0 || tokens < 0 {
		return errors.New("invalid quota")
	}
	q.mu.Lock()
	defer q.mu.Unlock()

	u, err := q.Store.Get(ctx, scope)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if u.ResetAt.IsZero() || !now.Before(u.ResetAt) {
		u = Usage{ResetAt: now.Add(24 * time.Hour)}
	}
	if u.Requests+1 > plan.RequestsPerDay || u.Tokens+tokens > plan.TokensPerDay {
		return errors.New("AI usage quota exceeded")
	}
	u.Requests++
	u.Tokens += tokens
	return q.Store.Put(ctx, scope, u)
}
