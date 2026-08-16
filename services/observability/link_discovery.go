package observability

// LinkObservation represents measured capacity and health for an FTN inter-node link.
type LinkObservation struct {
	LinkID string
	CapacityMbps float64
	LatencyMS float64
	LossPercent float64
	Healthy bool
}

func (o LinkObservation) Valid() bool {
	return o.LinkID != "" && o.CapacityMbps > 0 && o.LatencyMS >= 0 && o.LossPercent >= 0
}
