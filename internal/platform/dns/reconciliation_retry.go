package dns

import (
	"context"
	"fmt"
	"math"
	"time"
)

type RetryPolicy struct {
	MaxAttempts uint32 `json:"max_attempts"`
	BaseDelay time.Duration `json:"base_delay"`
	MaxDelay time.Duration `json:"max_delay"`
}

func (p RetryPolicy) Validate() error {
	if p.MaxAttempts == 0 { return fmt.Errorf("max attempts must be positive") }
	if p.BaseDelay <= 0 || p.MaxDelay <= 0 { return fmt.Errorf("retry delays must be positive") }
	if p.BaseDelay > p.MaxDelay { return fmt.Errorf("base delay cannot exceed max delay") }
	return nil
}

func RetryDelay(policy RetryPolicy, attempt uint32) time.Duration {
	if attempt == 0 { return 0 }
	factor := math.Pow(2, float64(attempt-1))
	d := time.Duration(float64(policy.BaseDelay) * factor)
	if d > policy.MaxDelay || d < 0 { return policy.MaxDelay }
	return d
}

// ExecuteWithRetry retries an already-approved reconciliation. The executor
// must remain idempotent; the caller should persist the reconciliation key.
func ExecuteWithRetry(ctx context.Context, policy RetryPolicy, execute func(context.Context) error) error {
	if err := policy.Validate(); err != nil { return err }
	if execute == nil { return fmt.Errorf("execute callback is required") }
	var lastErr error
	for attempt := uint32(1); attempt <= policy.MaxAttempts; attempt++ {
		if err := execute(ctx); err == nil { return nil } else { lastErr = err }
		if attempt == policy.MaxAttempts { break }
		delay := RetryDelay(policy, attempt)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("reconciliation failed after %d attempts: %w", policy.MaxAttempts, lastErr)
}
