package observability

// LinkQoSBudget describes the traffic budget for one FTN inter-node link.
type LinkQoSBudget struct {
	LinkID string
	CapacityMbps float64
	LiveReservedMbps float64
	ControlReservedMbps float64
	MigrationMaxMbps float64
}

func (b LinkQoSBudget) Valid() bool {
	return b.LinkID != "" && b.CapacityMbps > 0 && b.LiveReservedMbps >= 0 && b.ControlReservedMbps >= 0 && b.MigrationMaxMbps >= 0 && b.LiveReservedMbps+b.ControlReservedMbps+b.MigrationMaxMbps <= b.CapacityMbps
}
