package mesh

import "math"

// PathCandidate represents a route produced from the current mesh state.
type PathCandidate struct {
	Nodes []string `json:"nodes"`
	Cost uint64 `json:"cost"`
}

// ScoreLink combines routing metric with observed latency and packet loss.
// Lower scores are preferred. The caller decides policy/approval before
// applying any resulting route to a dataplane.
func ScoreLink(metric uint32, latencyMS, lossPercent float64) uint64 {
	if latencyMS < 0 { latencyMS = 0 }
	if lossPercent < 0 { lossPercent = 0 }
	v := float64(metric) + latencyMS*10 + lossPercent*100
	if math.IsInf(v, 0) || v >= float64(^uint64(0)) { return ^uint64(0) }
	return uint64(v)
}
