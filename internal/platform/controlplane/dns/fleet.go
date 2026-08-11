package dns

import (
	"sort"
	"time"
)

// Fleet ranks currently available DNS servers without changing traffic or
// configuration. Execution remains behind the Control Plane approval policy.
type Fleet struct { Servers []ServerHealth `json:"servers"` }

type RankedServer struct {
	ServerHealth
	Score float64 `json:"score"`
}

func (f Fleet) Rank(now time.Time, maxAge time.Duration) []RankedServer {
	out := make([]RankedServer, 0, len(f.Servers))
	for _, s := range f.Servers {
		if !Available(s, now, maxAge) { continue }
		latencyScore := 100.0 - min(s.LatencyMs, 100.0)
		servfailScore := 100.0 - s.ServfailRate
		score := latencyScore*0.55 + servfailScore*0.25 + min(s.QPS/100.0, 100.0)*0.20
		out = append(out, RankedServer{ServerHealth: s, Score: score})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func min(a, b float64) float64 { if a < b { return a }; return b }
