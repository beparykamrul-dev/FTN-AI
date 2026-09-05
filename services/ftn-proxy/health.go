package proxy

import "time"

type HealthState uint8

const (
	HealthUnknown HealthState = iota
	HealthHealthy
	HealthUnhealthy
)

type HealthCheck struct {
	Interval     time.Duration
	Timeout      time.Duration
	FailureLimit int
	SuccessLimit int
}

func DefaultHealthCheck() HealthCheck {
	return HealthCheck{Interval: 10 * time.Second, Timeout: 3 * time.Second, FailureLimit: 3, SuccessLimit: 2}
}

func (h HealthCheck) Valid() bool {
	return h.Interval > 0 && h.Timeout > 0 && h.Timeout <= h.Interval && h.FailureLimit > 0 && h.SuccessLimit > 0
}

type HealthTracker struct {
	State     HealthState
	Failures  int
	Successes int
}

func (h *HealthTracker) Observe(success bool, policy HealthCheck) HealthState {
	if h == nil {
		return HealthUnknown
	}
	if !policy.Valid() {
		policy = DefaultHealthCheck()
	}
	if success {
		h.Failures = 0
		if h.Successes < int(^uint(0)>>1) {
			h.Successes++
		}
		if h.Successes >= policy.SuccessLimit {
			h.State = HealthHealthy
		}
		return h.State
	}
	h.Successes = 0
	if h.Failures < int(^uint(0)>>1) {
		h.Failures++
	}
	if h.Failures >= policy.FailureLimit {
		h.State = HealthUnhealthy
	}
	return h.State
}
