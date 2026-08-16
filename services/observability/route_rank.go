package observability

// RouteRank combines backend health and capacity for best-fit routing.
type RouteRank struct {
	HealthScore float64
	CapacityScore float64
	PreferLocal bool
}

func (r RouteRank) Score() float64 {
	localBonus := 0.0
	if r.PreferLocal { localBonus = 5 }
	return r.HealthScore*0.6 + r.CapacityScore*0.4 + localBonus
}
