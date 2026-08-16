package observability

// QoSRecalculator derives a new link budget from a measured link.
type QoSRecalculator struct{}

func (QoSRecalculator) Recalculate(o LinkObservation) LinkQoSBudget {
	if !o.Valid() || !o.Healthy { return LinkQoSBudget{} }
	migration := DynamicQoS{}.MigrationBudget(o.CapacityMbps, o.LatencyMS, o.LossPercent)
	live := o.CapacityMbps * 0.60
	control := o.CapacityMbps * 0.10
	if live+control+migration > o.CapacityMbps {
		migration = o.CapacityMbps - live - control
		if migration < 0 { migration = 0 }
	}
	return LinkQoSBudget{LinkID:o.LinkID, CapacityMbps:o.CapacityMbps, LiveReservedMbps:live, ControlReservedMbps:control, MigrationMaxMbps:migration}
}
