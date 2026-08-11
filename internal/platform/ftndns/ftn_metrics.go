package ftndns

import "sync/atomic"

// FTNMetrics is the FTNDNS-native metrics contract. Exporters can map these
// counters to Prometheus/OpenTelemetry later without coupling the resolver to
// a particular telemetry vendor.
type FTNMetrics struct {
	Queries atomic.Uint64
	CacheHits atomic.Uint64
	CacheMisses atomic.Uint64
	NegativeHits atomic.Uint64
	UpstreamErrors atomic.Uint64
	Prefetches atomic.Uint64
	Deduplicated atomic.Uint64
	StampedePrevented atomic.Uint64
	MeshPeerHits atomic.Uint64
}

func (m *FTNMetrics) Snapshot() map[string]uint64 {
	return map[string]uint64{
		"queries":m.Queries.Load(), "cache_hits":m.CacheHits.Load(),
		"cache_misses":m.CacheMisses.Load(), "negative_hits":m.NegativeHits.Load(),
		"upstream_errors":m.UpstreamErrors.Load(), "prefetches":m.Prefetches.Load(),
		"deduplicated":m.Deduplicated.Load(), "stampede_prevented":m.StampedePrevented.Load(),
		"mesh_peer_hits":m.MeshPeerHits.Load(),
	}
}
