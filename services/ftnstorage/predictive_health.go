package ftnstorage

// PredictiveHealth provides a bounded, explainable risk indicator from recent
// storage observations. It does not mutate storage.
type PredictiveHealth struct {
	RecentFailures uint32
	ScrubErrors    uint32
	CapacityRatio  float64
}

func (p PredictiveHealth) Risk() float64 {
	risk := float64(p.RecentFailures)*0.2 + float64(p.ScrubErrors)*0.4
	if p.CapacityRatio > 0.8 {
		risk += (p.CapacityRatio - 0.8) * 2
	}
	if risk > 1 { return 1 }
	if risk < 0 { return 0 }
	return risk
}
