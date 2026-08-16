package observability

// RebalanceGuard prevents unsafe replica movement during unhealthy destination states.
type RebalanceGuard struct {
	MinTargetFreeGB float64
	RequireHealthyTarget bool
	RequireDifferentDomain bool
}

func (g RebalanceGuard) Allows(source, target ReplicaTarget) bool {
	if source.NodeID == "" || !target.Eligible() { return false }
	if g.RequireHealthyTarget && !target.Healthy { return false }
	if target.FreeGB < g.MinTargetFreeGB { return false }
	if g.RequireDifferentDomain && source.Domain.ID != "" && source.Domain.ID == target.Domain.ID { return false }
	return source.NodeID != target.NodeID
}
