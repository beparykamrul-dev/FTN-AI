package telemetry

import (
	"math"
	"strings"
	"time"
)

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
	return strings.TrimSpace(h.NodeID) != "" && !h.ObservedAt.IsZero() && finiteRange(h.CPUPercent, 0, 100) && finiteRange(h.MemoryPercent, 0, 100) && finiteRange(h.NetworkMbps, 0, math.MaxFloat64) && finiteRange(h.DNSQPS, 0, math.MaxFloat64)
}

func finiteRange(v, min, max float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= min && v <= max
}

func Fresh(h Heartbeat, now time.Time, maxAge time.Duration) bool {
	if !h.Valid() || maxAge < 0 || now.IsZero() { return false }
	now = now.UTC()
	observed := h.ObservedAt.UTC()
	if observed.After(now) { return false }
	return now.Sub(observed) <= maxAge
}
