package observability

// SignalPriority orders telemetry preservation when a backend is saturated.
type SignalPriority uint8

const (
	PriorityDebug SignalPriority = iota
	PriorityMetrics
	PriorityTraces
	PriorityAlerts
	PrioritySecurityAudit
)

// BackpressurePolicy defines the minimum priority accepted under pressure.
type BackpressurePolicy struct {
	MinPriority SignalPriority
}

func (p BackpressurePolicy) Allows(priority SignalPriority) bool {
	return priority >= p.MinPriority
}
