package ftndns

import (
	"sync"
	"time"
)

type NegativeCache struct {
	mu sync.RWMutex
	entries map[string]time.Time
	maxTTL time.Duration
}

func NewNegativeCache(maxTTL time.Duration) *NegativeCache {
	if maxTTL <= 0 { maxTTL = 30*time.Second }
	return &NegativeCache{entries: make(map[string]time.Time), maxTTL: maxTTL}
}

func (c *NegativeCache) Set(key string, ttl time.Duration, now time.Time) {
	if ttl <= 0 || ttl > c.maxTTL { ttl = c.maxTTL }
	c.mu.Lock(); c.entries[key] = now.Add(ttl); c.mu.Unlock()
}

func (c *NegativeCache) Hit(key string, now time.Time) bool {
	c.mu.RLock(); expires, ok := c.entries[key]; c.mu.RUnlock()
	if !ok { return false }
	if !expires.After(now) { c.mu.Lock(); delete(c.entries,key); c.mu.Unlock(); return false }
	return true
}

func (c *NegativeCache) Purge(now time.Time) int {
	c.mu.Lock(); defer c.mu.Unlock(); removed := 0
	for key, expires := range c.entries { if !expires.After(now) { delete(c.entries,key); removed++ } }
	return removed
}
