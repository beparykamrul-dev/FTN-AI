package proxy

import (
	"sort"
	"time"
)

// Upstream is the normalized FTN proxy origin state.
type Upstream struct {
	ID        string
	Address   string
	Healthy   bool
	Secure    bool
	LatencyMS int64
	Load      float64
	LastCheck time.Time
}

// RankUpstreams selects only usable origins and deterministically prefers
// secure, healthy, low-latency and low-load targets.
func RankUpstreams(items []Upstream) []Upstream {
	out := make([]Upstream, 0, len(items))
	for _, u := range items {
		if u.Healthy && u.Secure {
			out = append(out, u)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		if out[i].Load != out[j].Load { return out[i].Load < out[j].Load }
		return out[i].ID < out[j].ID
	})
	return out
}
