package observability

// TokenBucket is a bounded rate-control model for background FTN traffic.
type TokenBucket struct {
	RateMbps float64
	BurstMB float64
}

func (b TokenBucket) Valid() bool { return b.RateMbps > 0 && b.BurstMB > 0 }
