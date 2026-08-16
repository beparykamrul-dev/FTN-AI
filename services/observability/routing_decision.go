package observability

// RoutingDecision records why a telemetry backend was selected or rejected.
type RoutingDecision struct {
	Signal string
	Backend string
	Score float64
	Allowed bool
	Reason string
}

func (d RoutingDecision) Valid() bool {
	return d.Signal != "" && d.Backend != "" && d.Reason != ""
}
