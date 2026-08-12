package transport

import (
	"net"
	"sync"
	"time"
)

type IPAdmission struct {
	Limiter *IPRateLimiter
	mu      sync.Mutex
	blocked map[string]time.Time
	BlockFor time.Duration
}

func NewIPAdmission(limiter *IPRateLimiter, blockFor time.Duration) *IPAdmission {
	if blockFor <= 0 { blockFor = 10 * time.Minute }
	return &IPAdmission{Limiter: limiter, blocked: make(map[string]time.Time), BlockFor: blockFor}
}

func (a *IPAdmission) Allow(ip net.IP, now time.Time) bool {
	if a == nil || a.Limiter == nil { return false }
	key := ip.String()
	a.mu.Lock()
	until, exists := a.blocked[key]
	if exists && now.Before(until) { a.mu.Unlock(); return false }
	if exists { delete(a.blocked, key) }
	a.mu.Unlock()
	return a.Limiter.Allow(ip, now)
}

func (a *IPAdmission) Block(ip net.IP, now time.Time) {
	if a == nil { return }
	a.mu.Lock()
	a.blocked[ip.String()] = now.Add(a.BlockFor)
	a.mu.Unlock()
}

func (a *IPAdmission) Cleanup(now time.Time) {
	if a == nil { return }
	a.mu.Lock()
	defer a.mu.Unlock()
	for ip, until := range a.blocked { if !now.Before(until) { delete(a.blocked, ip) } }
	a.Limiter.Cleanup(now)
}
