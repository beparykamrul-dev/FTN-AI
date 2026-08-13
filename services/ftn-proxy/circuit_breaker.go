package proxy

import "time"

// CircuitState protects the FTN proxy from repeatedly sending traffic to a
// failing upstream.
type CircuitState uint8

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	Failures      int
	Threshold     int
	OpenedAt      time.Time
	Cooldown      time.Duration
}

func (c *CircuitBreaker) State(now time.Time) CircuitState {
	threshold := c.Threshold
	if threshold <= 0 { threshold = 5 }
	if c.Failures < threshold { return CircuitClosed }
	cooldown := c.Cooldown
	if cooldown <= 0 { cooldown = 30 * time.Second }
	if now.Sub(c.OpenedAt) >= cooldown { return CircuitHalfOpen }
	return CircuitOpen
}

func (c *CircuitBreaker) Success() {
	c.Failures = 0
	c.OpenedAt = time.Time{}
}

func (c *CircuitBreaker) Failure(now time.Time) {
	c.Failures++
	threshold := c.Threshold
	if threshold <= 0 { threshold = 5 }
	if c.Failures >= threshold && c.OpenedAt.IsZero() { c.OpenedAt = now }
}
