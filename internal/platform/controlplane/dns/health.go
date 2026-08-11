package dns

import "time"

type ServerHealth struct {
	NodeID string `json:"node_id"`
	Resolver string `json:"resolver"`
	Healthy bool `json:"healthy"`
	LatencyMs float64 `json:"latency_ms"`
	QPS float64 `json:"qps"`
	ServfailRate float64 `json:"servfail_rate"`
	ObservedAt time.Time `json:"observed_at"`
}

func (h ServerHealth) Valid() bool {
	return h.NodeID != "" && h.Resolver != "" && !h.ObservedAt.IsZero() && h.LatencyMs >= 0 && h.QPS >= 0 && h.ServfailRate >= 0 && h.ServfailRate <= 100
}

func Available(h ServerHealth, now time.Time, maxAge time.Duration) bool {
	return h.Valid() && maxAge >= 0 && now.Sub(h.ObservedAt) <= maxAge && h.Healthy
}
