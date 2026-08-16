package observability

// TelemetryAdmission decides whether a signal should enter the pipeline under pressure.
type TelemetryAdmission struct {
	Priority SignalPriority
	Policy BackpressurePolicy
}

func (a TelemetryAdmission) Allowed() bool {
	return a.Policy.Allows(a.Priority)
}
