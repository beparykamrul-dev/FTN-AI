package observability

// RoutingPolicy controls eligibility for telemetry backend selection.
type RoutingPolicy struct {
	MinScore float64
	PreferLocal bool
}

func (p RoutingPolicy) Allows(h BackendHealth) bool {
	return h.Healthy && h.Score() >= p.MinScore
}
