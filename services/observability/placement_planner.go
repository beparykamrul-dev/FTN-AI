package observability

// PlacementPlanner selects eligible replica targets while preserving failure-domain diversity.
type PlacementPlanner struct{}

func (PlacementPlanner) Plan(targets []ReplicaTarget, replicas uint32) []ReplicaTarget {
	return SelectReplicaTargets(targets, replicas)
}
