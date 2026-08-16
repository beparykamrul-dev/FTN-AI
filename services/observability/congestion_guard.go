package observability

// CongestionGuard prevents background migration from consuming unsafe network capacity.
type CongestionGuard struct {
	MaxLatencyMS float64
	MaxLossPercent float64
}

func (g CongestionGuard) Allows(latencyMS, lossPercent float64) bool {
	return latencyMS >= 0 && lossPercent >= 0 && latencyMS <= g.MaxLatencyMS && lossPercent <= g.MaxLossPercent
}
