package ftndns

import (
	"strings"
	"sync"
)

// CacheAdmission tracks query frequency so the resolver can avoid admitting
// one-off keys into a bounded hot cache.
type CacheAdmission struct { mu sync.Mutex; hits map[string]uint64; threshold uint64 }

func NewCacheAdmission(threshold uint64) *CacheAdmission { if threshold < 1 { threshold = 2 }; return &CacheAdmission{hits:make(map[string]uint64),threshold:threshold} }
func (a *CacheAdmission) Observe(key string) bool { if a == nil { return false }; key = strings.TrimSpace(key); if key == "" { return false }; a.mu.Lock(); defer a.mu.Unlock(); if a.hits == nil { a.hits = make(map[string]uint64) }; if a.hits[key] < a.threshold { a.hits[key]++ }; return a.hits[key] >= a.threshold }
func (a *CacheAdmission) Reset(key string) { if a == nil { return }; key = strings.TrimSpace(key); if key == "" { return }; a.mu.Lock(); delete(a.hits,key); a.mu.Unlock() }
