package provider

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("provider circuit open")
var ErrDuplicateOperation = errors.New("duplicate operation")

type CircuitState uint8
const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	mu sync.Mutex
	state CircuitState
	failures int
	threshold int
	openedAt time.Time
	cooldown time.Duration
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	if threshold < 1 { threshold = 1 }
	if cooldown <= 0 { cooldown = 5 * time.Second }
	return &CircuitBreaker{threshold: threshold, cooldown: cooldown}
}

func (c *CircuitBreaker) Allow(now time.Time) bool {
	c.mu.Lock(); defer c.mu.Unlock()
	if c.state == CircuitClosed { return true }
	if c.state == CircuitOpen && now.Sub(c.openedAt) >= c.cooldown {
		c.state = CircuitHalfOpen
		return true
	}
	return c.state == CircuitHalfOpen
}

func (c *CircuitBreaker) Success() {
	c.mu.Lock(); defer c.mu.Unlock()
	c.failures = 0
	c.state = CircuitClosed
}

func (c *CircuitBreaker) Failure(now time.Time) {
	c.mu.Lock(); defer c.mu.Unlock()
	c.failures++
	if c.state == CircuitHalfOpen || c.failures >= c.threshold {
		c.state = CircuitOpen
		c.openedAt = now
	}
}

type IdempotencyRegistry struct {
	mu sync.Mutex
	keys map[string]time.Time
}

func NewIdempotencyRegistry() *IdempotencyRegistry { return &IdempotencyRegistry{keys: make(map[string]time.Time)} }

func (r *IdempotencyRegistry) Claim(key string, now time.Time, ttl time.Duration) bool {
	if key == "" { return true }
	r.mu.Lock(); defer r.mu.Unlock()
	if expires, ok := r.keys[key]; ok && now.Before(expires) { return false }
	r.keys[key] = now.Add(ttl)
	return true
}
