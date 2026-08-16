package observability

// AlertRule defines a normalized threshold rule for FTN telemetry.
type AlertRule struct {
	Name string
	Service string
	Metric string
	Operator string
	Threshold float64
}

func (r AlertRule) Valid() bool {
	return r.Name != "" && r.Service != "" && r.Metric != "" && r.Operator != ""
}
