package runtime

import (
	"context"
	"time"
)

type HealthState string

const (
	HealthHealthy HealthState = "healthy"
	HealthFailed  HealthState = "failed"
	HealthUnknown HealthState = "unknown"
)

type HealthResult struct {
	Service   string
	State     HealthState
	CheckedAt time.Time
	Error     string
}

func CheckHealth(ctx context.Context, manager *ServiceManager) []HealthResult {
	health := manager.Health(ctx)
	results := make([]HealthResult, 0, len(health))
	now := time.Now().UTC()

	for name, err := range health {
		result := HealthResult{
			Service:   name,
			State:     HealthHealthy,
			CheckedAt: now,
		}
		if err != nil {
			result.State = HealthFailed
			result.Error = err.Error()
		}
		results = append(results, result)
	}
	return results
}
