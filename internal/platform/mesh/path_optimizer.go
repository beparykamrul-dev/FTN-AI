package mesh

import "math"

// PathCandidate is the shared route candidate model used by mesh ranking and
// cost-based optimization.
type PathCandidate struct {
	Nodes []string `json:"nodes,omitempty"`
	Cost uint64 `json:"cost,omitempty"`
	PeerID string `json:"peer_id,omitempty"`
	LatencyMS float64 `json:"latency_ms,omitempty"`
	LossPct float64 `json:"loss_pct,omitempty"`
	CapacityMbps float64 `json:"capacity_mbps,omitempty"`
	HealthScore uint8 `json:"health_score,omitempty"`
}

func ScoreLink(metric uint32, latencyMS, lossPercent float64) uint64 {
	if latencyMS < 0 { latencyMS = 0 }
	if lossPercent < 0 { lossPercent = 0 }
	v := float64(metric) + latencyMS*10 + lossPercent*100
	if math.IsInf(v, 0) || v >= float64(^uint64(0)) { return ^uint64(0) }
	return uint64(v)
}
