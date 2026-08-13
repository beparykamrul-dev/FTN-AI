package ftnwan

// PathMetrics describes measured network health for a candidate path.
type PathMetrics struct {
	Name        string
	LatencyMs   float64
	JitterMs    float64
	LossPct     float64
	CapacityPct float64
	Healthy     bool
}

// Score produces a deterministic health score. Higher is better.
func Score(m PathMetrics) float64 {
	if !m.Healthy {
		return -1
	}
	latency := 100.0 / (1.0 + m.LatencyMs)
	jitter := 100.0 / (1.0 + m.JitterMs)
	loss := 100.0 - m.LossPct*10.0
	capacity := 100.0 - m.CapacityPct
	if loss < 0 { loss = 0 }
	if capacity < 0 { capacity = 0 }
	return latency*0.35 + jitter*0.20 + loss*0.25 + capacity*0.20
}

// SelectBest returns the highest-scoring healthy path without applying any
// network change. A separate policy/deployment layer must approve changes.
func SelectBest(paths []PathMetrics) (PathMetrics, bool) {
	var best PathMetrics
	bestScore := -1.0
	found := false
	for _, p := range paths {
		s := Score(p)
		if s > bestScore {
			best, bestScore, found = p, s, true
		}
	}
	return best, found
}
