package telemetry

import "time"

type Heartbeat struct {
	NodeID string `json:"node_id"`
	ObservedAt time.Time `json:"observed_at"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkMbps float64 `json:"network_mbps"`
	DNSQPS float64 `json:"dns_qps"`
	Healthy bool `json:"healthy"`
}

func (h Heartbeat) Valid() bool {
	return h.NodeID != "" && !h.ObservedAt.IsZero() && h.CPUPercent >= 0 && h.CPUPercent <= 100 && h.MemoryPercent >= 0 && h.MemoryPercent <= 100 && h.NetworkMbps >= 0 && h.DNSQPS >= 0
}

func Fresh(h Heartbeat, now time.Time, maxAge time.Duration) bool {
	if !h.Valid() || maxAge < 0 { return false }
	return now.Sub(h.ObservedAt) <= maxAge
}
