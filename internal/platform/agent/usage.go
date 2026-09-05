package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// UsageGate enforces the canonical plan entitlements defined in quota.go.
type UsageGate struct {
	mu    sync.Mutex
	plans map[string]Plan
	usage map[string]Usage
	now   func() time.Time
}

func NewUsageGate(plans []Plan) *UsageGate {
	m := make(map[string]Plan, len(plans))
	for _, p := range plans {
		if p.ID != "" {
			m[p.ID] = p
		}
	}
	return &UsageGate{plans: m, usage: make(map[string]Usage), now: time.Now}
}

func (g *UsageGate) CheckAndConsume(ctx context.Context, principal, planID string, tokens int64) error {
	if g == nil {
		return fmt.Errorf("usage gate is required")
	}
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || planID == "" || tokens < 0 {
		return fmt.Errorf("invalid usage request")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.plans[planID]
	if !ok {
		return fmt.Errorf("unknown plan: %s", planID)
	}
	if p.RequestsPerDay <= 0 || p.TokensPerDay <= 0 {
		return fmt.Errorf("invalid plan: %s", planID)
	}
	now := g.now
	if now == nil {
		now = time.Now
	}
	current := now()
	u := g.usage[principal]
	if u.ResetAt.IsZero() || !current.Before(u.ResetAt) {
		u = Usage{ResetAt: current.UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)}
	}
	if u.Requests >= p.RequestsPerDay {
		return fmt.Errorf("request quota exceeded")
	}
	if u.Tokens+tokens > p.TokensPerDay {
		return fmt.Errorf("token quota exceeded")
	}
	u.Requests++
	u.Tokens += tokens
	g.usage[principal] = u
	return nil
}
