package observability

// StoragePath represents a candidate multi-hop path for FTN replica transfer.
type StoragePath struct {
	PathID string
	Nodes []string
	LatencyMS float64
	LossPercent float64
	CapacityMbps float64
	Healthy bool
}

func (p StoragePath) Eligible() bool {
	return p.PathID != "" && len(p.Nodes) >= 2 && p.LatencyMS >= 0 && p.LossPercent >= 0 && p.CapacityMbps > 0 && p.Healthy
}
