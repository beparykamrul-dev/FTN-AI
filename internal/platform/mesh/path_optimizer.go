package mesh

import "math"

// PathCandidate is the canonical route candidate model used by both optimizer
// and path-selection layers.
type PathCandidate struct {
	Nodes []string `json:"nodes,omitempty"`
	PeerID string `json:"peer_id,omitempty"`
	Cost uint64 `json:"cost,omitempty"`
	LatencyMS float64 `json:"latency_ms,omitempty"`
	LossPct float64 `json:"loss_pct,omitempty"`
	CapacityMbps float64 `json:"capacity_mbps,omitempty"`
	HealthScore uint8 `json:"health_score,omitempty"`
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
