package observability

// WeightedPathScheduler ranks healthy paths by usable capacity and network quality.
type WeightedPathScheduler struct{}

func (WeightedPathScheduler) Score(p StoragePath) float64 {
	if !p.Eligible() { return 0 }
	latencyFactor := 1.0 / (1.0 + p.LatencyMS)
	lossFactor := 1.0 / (1.0 + p.LossPercent)
	return p.CapacityMbps * latencyFactor * lossFactor
}

func (s WeightedPathScheduler) Rank(paths []StoragePath) []StoragePath {
	result := make([]StoragePath, 0, len(paths))
	for _, p := range paths { if p.Eligible() { result = append(result, p) } }
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if s.Score(result[j]) > s.Score(result[i]) { result[i], result[j] = result[j], result[i] }
		}
	}
	return result
}
