package kernel

import (
	"math"
	"sort"
	"sync"
	"time"
)

type GlobalPathSample struct { Source string; Destination string; RTT time.Duration; Jitter time.Duration; Loss float64; Utilization float64; Healthy bool; ObservedAt time.Time }
type PathScore struct { Sample GlobalPathSample; Score float64 }

type GlobalPathRouter struct {
	mu sync.RWMutex
	samples map[string][]GlobalPathSample
	maxAge time.Duration
}

func NewGlobalPathRouter(maxAge time.Duration) *GlobalPathRouter {
	if maxAge <= 0 { maxAge = 15 * time.Second }
	return &GlobalPathRouter{samples: make(map[string][]GlobalPathSample), maxAge: maxAge}
}
func pathKey(source, destination string) string { return source+"->"+destination }
func (r *GlobalPathRouter) Observe(s GlobalPathSample) {
	if s.Source=="" || s.Destination=="" || s.RTT<0 || s.Jitter<0 || s.Loss<0 || s.Loss>1 || s.Utilization<0 || s.Utilization>1 { return }
	r.mu.Lock(); defer r.mu.Unlock()
	key := pathKey(s.Source,s.Destination)
	r.samples[key] = append(r.samples[key], s)
	// Retain a bounded recent history; candidates are filtered by maxAge below.
	if len(r.samples[key]) > 64 { r.samples[key] = r.samples[key][len(r.samples[key])-64:] }
}
func (r *GlobalPathRouter) Candidates(source,destination string,now time.Time) []PathScore {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]PathScore,0)
	for _, s := range r.samples[pathKey(source,destination)] {
		if !s.Healthy || (!s.ObservedAt.IsZero() && now.Sub(s.ObservedAt)>r.maxAge) { continue }
		out=append(out,PathScore{Sample:s,Score:scorePath(s)})
	}
	sort.Slice(out,func(i,j int)bool{return out[i].Score<out[j].Score})
	return out
}
func scorePath(s GlobalPathSample) float64 { return float64(s.RTT.Milliseconds())+0.35*float64(s.Jitter.Milliseconds())+1000*s.Loss+200*s.Utilization }
func IsUsable(s GlobalPathSample) bool { return s.Healthy&&s.RTT>=0&&s.Loss>=0&&s.Loss<=1&&s.Utilization>=0&&s.Utilization<=1&&!math.IsNaN(float64(s.RTT)) }
