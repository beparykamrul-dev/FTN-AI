package observability

// AdaptiveBandwidth computes a bounded background migration rate from network conditions.
type AdaptiveBandwidth struct {
	MinMbps float64
	MaxMbps float64
	TargetLatencyMS float64
	TargetLossPercent float64
}

func (a AdaptiveBandwidth) Valid() bool {
	return a.MinMbps >= 0 && a.MaxMbps > a.MinMbps && a.TargetLatencyMS > 0 && a.TargetLossPercent >= 0
}

func (a AdaptiveBandwidth) Rate(latencyMS, lossPercent float64) float64 {
	if !a.Valid() { return 0 }
	if latencyMS > a.TargetLatencyMS*1.5 || lossPercent > a.TargetLossPercent*1.5 { return a.MinMbps }
	if latencyMS > a.TargetLatencyMS || lossPercent > a.TargetLossPercent { return (a.MinMbps+a.MaxMbps)/2 }
	return a.MaxMbps
}
