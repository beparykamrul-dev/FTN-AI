package observability

// MigrationBandwidthPolicy bounds background storage-transfer bandwidth.
type MigrationBandwidthPolicy struct {
	MaxMbps float64
	MinMbps float64
}

func (p MigrationBandwidthPolicy) Valid() bool {
	return p.MaxMbps > 0 && p.MinMbps >= 0 && p.MinMbps <= p.MaxMbps
}

func (p MigrationBandwidthPolicy) Clamp(requested float64) float64 {
	if requested < p.MinMbps { return p.MinMbps }
	if requested > p.MaxMbps { return p.MaxMbps }
	return requested
}
