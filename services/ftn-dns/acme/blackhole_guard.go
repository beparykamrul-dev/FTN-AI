package acme

import "time"

// BlackholeGuard coordinates certificate rotation with FTN security state.
// Certificate rotation may be checked hourly, but an active certificate is
// never discarded merely because a security event occurs.
type BlackholeGuard struct {
	RotationInterval time.Duration
	RequireHealthy   bool
}

// RotationDecision is consumed by the FTN certificate controller.
type RotationDecision struct {
	Rotate bool
	Reason string
}

// DecideRotation keeps the currently valid certificate during an unhealthy or
// blackhole condition and permits rotation when the security path is healthy.
func (g BlackholeGuard) DecideRotation(now, lastRotation time.Time, pathHealthy bool, certExpiring bool) RotationDecision {
	interval := g.RotationInterval
	if interval <= 0 {
		interval = time.Hour
	}
	if now.Sub(lastRotation) < interval {
		return RotationDecision{Reason: "rotation interval not reached"}
	}
	if g.RequireHealthy && !pathHealthy {
		return RotationDecision{Reason: "FTN security path unhealthy; retain current certificate"}
	}
	if !certExpiring {
		return RotationDecision{Reason: "certificate is outside renewal window"}
	}
	return RotationDecision{Rotate: true, Reason: "certificate rotation permitted"}
}
