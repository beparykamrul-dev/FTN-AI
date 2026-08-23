package kernel

import (
	"sort"
	"sync"
	"time"
)

// EndpointSample is a small in-memory observation used for route selection.
// It contains no credentials or packet payloads.
type EndpointSample struct {
	Endpoint   string
	RTT        time.Duration
	Loss       float64
	Healthy    bool
	ObservedAt time.Time
}

// LatencyRouter selects the healthiest low-latency endpoint without changing
// the security or authorization boundary of the kernel/backend.
type LatencyRouter struct {
	mu      sync.RWMutex
	samples map[string]EndpointSample
}

func NewLatencyRouter() *LatencyRouter {
	return &LatencyRouter{samples: make(map[string]EndpointSample)}
}

func (r *LatencyRouter) Observe(s EndpointSample) {
	if s.Endpoint == "" || s.RTT < 0 || s.Loss < 0 || s.Loss > 1 {
		return
	}
	r.mu.Lock()
	r.samples[s.Endpoint] = s
	r.mu.Unlock()
}

func (r *LatencyRouter) Best() (EndpointSample, bool) {
	r.mu.RLock()
	items := make([]EndpointSample, 0, len(r.samples))
	for _, s := range r.samples {
		if s.Healthy {
			items = append(items, s)
		}
	}
	r.mu.RUnlock()
	if len(items) == 0 {
		return EndpointSample{}, false
	}
	sort.Slice(items, func(i, j int) bool {
		// Prefer lower RTT; use packet loss as the tie breaker.
		if items[i].RTT != items[j].RTT {
			return items[i].RTT < items[j].RTT
		}
		return items[i].Loss < items[j].Loss
	})
	return items[0], true
}
