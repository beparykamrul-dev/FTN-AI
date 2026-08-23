package kernel

import (
	"math"
	"sort"
	"sync"
	"time"
)

// GlobalPathSample describes one measured path between FTN network points.
// It is intentionally independent of any particular transport protocol.
type GlobalPathSample struct {
	Source       string
	Destination  string
	RTT          time.Duration
	Jitter       time.Duration
	Loss         float64
	Utilization  float64
	Healthy      bool
	ObservedAt   time.Time
}

// PathScore is a normalized score where lower is better. The score combines
// latency, jitter, loss and congestion so DNS is not the only optimization
// target; application and backbone traffic use the same path-selection model.
type PathScore struct {
	Sample GlobalPathSample
	Score  float64
}

// GlobalPathRouter maintains recent path observations for the FTN backbone.
type GlobalPathRouter struct {
	mu       sync.RWMutex
	samples  map[string]GlobalPathSample
	maxAge   time.Duration
}

func NewGlobalPathRouter(maxAge time.Duration) *GlobalPathRouter {
	if maxAge <= 0 {
		maxAge = 15 * time.Second
	}
	return &GlobalPathRouter{samples: make(map[string]GlobalPathSample), maxAge: maxAge}
}

func pathKey(source, destination string) string { return source + "->" + destination }

func (r *GlobalPathRouter) Observe(s GlobalPathSample) {
	if s.Source == "" || s.Destination == "" || s.RTT < 0 || s.Jitter < 0 ||
		s.Loss < 0 || s.Loss > 1 || s.Utilization < 0 || s.Utilization > 1 {
		return
	}
	r.mu.Lock()
	r.samples[pathKey(s.Source, s.Destination)] = s
	r.mu.Unlock()
}

func (r *GlobalPathRouter) Candidates(source, destination string, now time.Time) []PathScore {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]PathScore, 0)
	for _, s := range r.samples {
		if !s.Healthy || s.Source != source || s.Destination != destination {
			continue
		}
		if !s.ObservedAt.IsZero() && now.Sub(s.ObservedAt) > r.maxAge {
			continue
		}
		out = append(out, PathScore{Sample: s, Score: scorePath(s)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score < out[j].Score })
	return out
}

func scorePath(s GlobalPathSample) float64 {
	// RTT is the dominant term. Jitter, loss and utilization prevent selecting
	// a path that is fast only when lightly loaded or unstable.
	rtt := float64(s.RTT.Milliseconds())
	jitter := float64(s.Jitter.Milliseconds())
	return rtt + 0.35*jitter + 1000*s.Loss + 200*s.Utilization
}

// IsUsable rejects obviously degraded paths before they enter the candidate set.
func IsUsable(s GlobalPathSample) bool {
	return s.Healthy && s.RTT >= 0 && s.Loss >= 0 && s.Loss <= 1 &&
		s.Utilization >= 0 && s.Utilization <= 1 && !math.IsNaN(float64(s.RTT))
}
