package ftndns

import "sync"

// CacheAdmission tracks query frequency so the resolver can avoid admitting
// one-off keys into a bounded hot cache.
type CacheAdmission struct { mu sync.Mutex; hits map[string]uint64; threshold uint64 }

func NewCacheAdmission(threshold uint64) *CacheAdmission { if threshold < 1 { threshold = 2 }; return &CacheAdmission{hits:make(map[string]uint64),threshold:threshold} }
func (a *CacheAdmission) Observe(key string) bool { a.mu.Lock(); defer a.mu.Unlock(); a.hits[key]++; return a.hits[key] >= a.threshold }
func (a *CacheAdmission) Reset(key string) { a.mu.Lock(); delete(a.hits,key); a.mu.Unlock() }
