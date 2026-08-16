package observability

// FailureDomain identifies an independent placement domain for telemetry replicas.
type FailureDomain struct {
	ID string
	Kind string
	Healthy bool
}

func (d FailureDomain) Eligible() bool { return d.ID != "" && d.Kind != "" && d.Healthy }
