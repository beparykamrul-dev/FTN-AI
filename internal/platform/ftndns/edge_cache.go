package ftndns

import "time"

// EdgeCachePolicy describes the provider-neutral policy needed to evolve
// FTNDNS cache nodes toward future FTN Edge/CDN roles.
type EdgeCachePolicy struct {
	Enabled bool
	DefaultTTL time.Duration
	StaleWindow time.Duration
	Prefetch bool
}

func DefaultEdgeCachePolicy() EdgeCachePolicy {
	return EdgeCachePolicy{Enabled:true, DefaultTTL:5*time.Minute, StaleWindow:30*time.Second, Prefetch:true}
}
