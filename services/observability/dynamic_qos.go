package observability

// DynamicQoS derives a conservative migration budget from current link conditions.
type DynamicQoS struct{}

func (DynamicQoS) MigrationBudget(capacity, latencyMS, lossPercent float64) float64 {
	if capacity <= 0 || latencyMS < 0 || lossPercent < 0 { return 0 }
	budget := capacity * 0.20
	if latencyMS > 50 { budget *= 0.5 }
	if lossPercent > 1 { budget *= 0.5 }
	if budget < 0 { return 0 }
	return budget
}
