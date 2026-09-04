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
	Source      string
	Destination string
	RTT         time.Duration
	Jitter      time.Duration
	Loss        float64
	Utilization float64
	Healthy     bool
	ObservedAt  time.Time
}

type PathScore struct {
	Sample GlobalPathSample
	Score  float64
}

type GlobalPathRouter struct {
	mu      sync.RWMutex
	samples map[string][]GlobalPathSample
	maxAge  time.Duration
}

func NewGlobalPathRouter(maxAge time.Duration) *GlobalPathRouter {
	if maxAge <= 0 {
		maxAge = 15 * time.Second
	}
	return &GlobalPathRouter{samples: make(map[string][]GlobalPathSample), maxAge: maxAge}
}

func pathKey(source, destination string) string { return source + "->" + destination }

func (r *GlobalPathRouter) Observe(s GlobalPathSample) {
	if !IsUsable(s) {
		return
	}
	r.mu.Lock()
	key := pathKey(s.Source, s.Destination)
	history := append(r.samples[key], s)
	if len(history) > 32 {
		history = history[len(history)-32:]
	}
	r.samples[key] = history
	r.mu.Unlock()
}

func (r *GlobalPathRouter) Candidates(source, destination string, now time.Time) []PathScore {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]PathScore, 0)
	for _, s := range r.samples[pathKey(source, destination)] {
		if !IsUsable(s) {
			continue
		}
		if !s.ObservedAt.IsZero() && now.Sub(s.ObservedAt) > r.maxAge {
			continue
		}
		out = append(out, PathScore{Sample: s, Score: scorePath(s)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score < out[j].Score
		}
		return pathKey(out[i].Sample.Source, out[i].Sample.Destination) < pathKey(out[j].Sample.Source, out[j].Sample.Destination)
	})
	return out
}

func scorePath(s GlobalPathSample) float64 {
	rtt := float64(s.RTT.Milliseconds())
	jitter := float64(s.Jitter.Milliseconds())
	return rtt + 0.35*jitter + 1000*s.Loss + 200*s.Utilization
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func IsUsable(s GlobalPathSample) bool {
	return s.Healthy && s.RTT >= 0 && s.Jitter >= 0 &&
		finite(s.Loss) && s.Loss >= 0 && s.Loss <= 1 &&
		finite(s.Utilization) && s.Utilization >= 0 && s.Utilization <= 1
}
