package ftndns

import "context"

// EdgeProvider is the provider-neutral contract for future FTN Edge/CDN nodes.
// DNS remains FTNDNS-owned; content delivery is an independent capability.
type EdgeProvider interface {
	ID() string
	Health(ctx context.Context) error
	Purge(ctx context.Context, keys []string) error
	Warm(ctx context.Context, keys []string) error
}

type EdgeNodeState struct {
	Provider string `json:"provider"`
	Healthy bool `json:"healthy"`
	LatencyMS int64 `json:"latency_ms"`
}
