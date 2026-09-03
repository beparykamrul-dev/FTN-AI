package mesh

import (
	"math"
	"sort"
)

type PathPolicy struct {
	LatencyWeight float64
	LossWeight float64
	CapacityWeight float64
	HealthWeight float64
}

func DefaultPathPolicy() PathPolicy { return PathPolicy{LatencyWeight: 0.40, LossWeight: 0.25, CapacityWeight: 0.15, HealthWeight: 0.20} }

func RankPaths(paths []PathCandidate, policy PathPolicy) []PathCandidate {
	out := append([]PathCandidate(nil), paths...)
	score := func(p PathCandidate) float64 {
		latency := 1.0 / (1.0 + math.Max(0, p.LatencyMS))
		loss := 1.0 / (1.0 + math.Max(0, p.LossPct))
		capacity := math.Log1p(math.Max(0, p.CapacityMbps))
		health := float64(p.HealthScore) / 100.0
		return policy.LatencyWeight*latency + policy.LossWeight*loss + policy.CapacityWeight*capacity + policy.HealthWeight*health
	}
	sort.SliceStable(out, func(i, j int) bool { return score(out[i]) > score(out[j]) })
	return out
}
