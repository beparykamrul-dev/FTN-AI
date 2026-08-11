package ftndns

import (
	"sync"
	"time"
)

type RefreshGuard struct { mu sync.Mutex; active map[string]time.Time }
func NewRefreshGuard() *RefreshGuard { return &RefreshGuard{active:make(map[string]time.Time)} }
func (g *RefreshGuard) TryAcquire(key string, now time.Time, hold time.Duration) bool { g.mu.Lock(); defer g.mu.Unlock(); if until,ok:=g.active[key]; ok && until.After(now) { return false }; g.active[key]=now.Add(hold); return true }
func (g *RefreshGuard) Release(key string) { g.mu.Lock(); delete(g.active,key); g.mu.Unlock() }
