package telemetry

import "time"

type HealthStatus string

const (
	HealthHealthy   HealthStatus = "healthy"
	HealthDegraded  HealthStatus = "degraded"
	HealthUnhealthy HealthStatus = "unhealthy"
	HealthStale     HealthStatus = "stale"
)

type HealthSnapshot struct {
	NodeID string `json:"node_id"`
	Status HealthStatus `json:"status"`
	Reason string `json:"reason,omitempty"`
	ObservedAt time.Time `json:"observed_at"`
}

func Evaluate(h Heartbeat, now time.Time, maxAge time.Duration) HealthSnapshot {
	if !h.Valid() {
		return HealthSnapshot{NodeID: h.NodeID, Status: HealthUnhealthy, Reason: "invalid heartbeat", ObservedAt: h.ObservedAt}
	}
	if !Fresh(h, now, maxAge) {
		return HealthSnapshot{NodeID: h.NodeID, Status: HealthStale, Reason: "heartbeat expired", ObservedAt: h.ObservedAt}
	}
	if !h.Healthy {
		return HealthSnapshot{NodeID: h.NodeID, Status: HealthUnhealthy, Reason: "node reported unhealthy", ObservedAt: h.ObservedAt}
	}
	if h.CPUPercent >= 90 || h.MemoryPercent >= 90 {
		return HealthSnapshot{NodeID: h.NodeID, Status: HealthDegraded, Reason: "resource pressure", ObservedAt: h.ObservedAt}
	}
	return HealthSnapshot{NodeID: h.NodeID, Status: HealthHealthy, ObservedAt: h.ObservedAt}
}
