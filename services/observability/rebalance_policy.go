package observability

// RebalancePolicy determines when a replica movement plan is eligible.
type RebalancePolicy struct {
	MaxLoadPercent float64
	MinFreeGB float64
	MaxLatencyMS float64
}

func (p RebalancePolicy) NeedsRebalance(cap BackendCapacity, latencyMS float64) bool {
	if cap.Name == "" { return false }
	return cap.LoadPercent > p.MaxLoadPercent || cap.StoragePressurePercent > 90 || cap.QueueDepth > 1000 || cap.IngestPerSecond < 0 || cap.Name == "" || latencyMS > p.MaxLatencyMS
}
