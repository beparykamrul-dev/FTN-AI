package observability

// MigrationCost estimates the relative cost of moving pending data.
type MigrationCost struct {
	Bytes uint64
	BandwidthMbps float64
	LatencyMS float64
	PathPenalty float64
}

func (c MigrationCost) Valid() bool {
	return c.Bytes > 0 && c.BandwidthMbps > 0 && c.LatencyMS >= 0 && c.PathPenalty >= 0
}

func (c MigrationCost) EstimatedSeconds() float64 {
	if !c.Valid() { return 0 }
	bytesPerSecond := c.BandwidthMbps * 1_000_000 / 8
	return float64(c.Bytes)/bytesPerSecond + c.LatencyMS/1000 + c.PathPenalty
}
