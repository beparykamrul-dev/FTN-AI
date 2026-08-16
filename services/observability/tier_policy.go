package observability

// TierPolicy controls promotion and demotion thresholds for buffered telemetry.
type TierPolicy struct {
	PromoteLatencyMS float64
	DemoteFreePercent float64
}

func (p TierPolicy) Valid() bool {
	return p.PromoteLatencyMS > 0 && p.DemoteFreePercent >= 0 && p.DemoteFreePercent <= 100
}
