package acme

import "time"

// RotationPolicy controls automatic certificate renewal checks. The scheduler
// runs hourly, but renewal is performed only when the certificate is inside
// the configured renewal window; it does not force unsafe hourly reissuance.
type RotationPolicy struct {
	CheckInterval   time.Duration
	RenewBefore     time.Duration
	FallbackGrace   time.Duration
}

func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{
		CheckInterval: time.Hour,
		RenewBefore:   30 * 24 * time.Hour,
		FallbackGrace: 24 * time.Hour,
	}
}

// ShouldRenew determines whether the certificate should enter the ACME
// renewal workflow based on its expiry time.
func (p RotationPolicy) ShouldRenew(now, notAfter time.Time) bool {
	window := p.RenewBefore
	if window <= 0 { window = 30 * 24 * time.Hour }
	return !notAfter.IsZero() && !now.Before(notAfter.Add(-window))
}
