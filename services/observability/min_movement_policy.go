package observability

// MinMovementPolicy prevents frequent pending-chunk reassignment for small path changes.
type MinMovementPolicy struct {
	CapacityChangePercent float64
	LatencyChangePercent float64
	LossChangePercent float64
}

func (p MinMovementPolicy) Valid() bool {
	return p.CapacityChangePercent >= 0 && p.LatencyChangePercent >= 0 && p.LossChangePercent >= 0
}

func (p MinMovementPolicy) ShouldMove(capacityChange, latencyChange, lossChange float64) bool {
	if !p.Valid() { return false }
	return capacityChange >= p.CapacityChangePercent || latencyChange >= p.LatencyChangePercent || lossChange >= p.LossChangePercent
}
