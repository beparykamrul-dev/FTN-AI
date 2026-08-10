package mesh

import "time"

type ProbeResult struct {
	LinkID string `json:"link_id"`
	Healthy bool `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
	LossPercent float64 `json:"loss_percent"`
	ObservedAt time.Time `json:"observed_at"`
}

type HealthThresholds struct {
	MaxLatencyMS float64
	MaxLossPercent float64
}

func EvaluateProbe(r ProbeResult, t HealthThresholds) LinkState {
	if !r.Healthy || r.LossPercent >= 100 { return LinkDown }
	if r.LatencyMS > t.MaxLatencyMS || r.LossPercent > t.MaxLossPercent { return LinkDegraded }
	return LinkUp
}

func ProbeAge(now, observed time.Time) time.Duration {
	if observed.IsZero() { return time.Duration(1<<63 - 1) }
	if now.Before(observed) { return 0 }
	return now.Sub(observed)
}
