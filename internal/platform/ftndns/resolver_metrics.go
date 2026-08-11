package ftndns

import "sync/atomic"

type ResolverMetrics struct { Queries atomic.Uint64; CacheHits atomic.Uint64; CacheMisses atomic.Uint64; NegativeHits atomic.Uint64; UpstreamErrors atomic.Uint64; Prefetches atomic.Uint64 }
func (m *ResolverMetrics) RecordQuery(hit bool) { m.Queries.Add(1); if hit { m.CacheHits.Add(1) } else { m.CacheMisses.Add(1) } }
func (m *ResolverMetrics) Snapshot() map[string]uint64 { return map[string]uint64{"queries":m.Queries.Load(),"cache_hits":m.CacheHits.Load(),"cache_misses":m.CacheMisses.Load(),"negative_hits":m.NegativeHits.Load(),"upstream_errors":m.UpstreamErrors.Load(),"prefetches":m.Prefetches.Load()} }
