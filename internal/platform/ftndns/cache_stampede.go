package ftndns

import (
	"strings"
	"sync"
	"time"
)

type RefreshGuard struct { mu sync.Mutex; active map[string]time.Time }
func NewRefreshGuard() *RefreshGuard { return &RefreshGuard{active:make(map[string]time.Time)} }
func (g *RefreshGuard) TryAcquire(key string, now time.Time, hold time.Duration) bool { if g == nil { return false }; key = strings.TrimSpace(key); if key == "" || hold <= 0 { return false }; if now.IsZero() { now = time.Now().UTC() } else { now = now.UTC() }; g.mu.Lock(); defer g.mu.Unlock(); if g.active == nil { g.active = make(map[string]time.Time) }; if until,ok:=g.active[key]; ok && until.After(now) { return false }; g.active[key]=now.Add(hold); return true }
func (g *RefreshGuard) Release(key string) { if g == nil { return }; key = strings.TrimSpace(key); if key == "" { return }; g.mu.Lock(); delete(g.active,key); g.mu.Unlock() }
