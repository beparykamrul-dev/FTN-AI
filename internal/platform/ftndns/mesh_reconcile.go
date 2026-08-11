package ftndns

import "sort"

type MeshNode struct {
	ID ProviderID
	Healthy bool
	LatencyMS int64
	Snapshot Snapshot
}

type MeshReconcilePlan struct {
	Reference ProviderID
	Targets []ProviderID
	Reason string
}

// BuildMeshReconcilePlan selects the healthiest, lowest-latency consistent
// reference. It only produces a plan; it never mutates provider state.
func BuildMeshReconcilePlan(nodes []MeshNode, consistency ConsistencyResult) MeshReconcilePlan {
	if len(nodes) == 0 || consistency.Consistent {
		return MeshReconcilePlan{}
	}
	candidates := append([]MeshNode(nil), nodes...)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Healthy != candidates[j].Healthy { return candidates[i].Healthy }
		if candidates[i].LatencyMS != candidates[j].LatencyMS { return candidates[i].LatencyMS < candidates[j].LatencyMS }
		return candidates[i].ID < candidates[j].ID
	})
	if !candidates[0].Healthy { return MeshReconcilePlan{Reason: "no healthy mesh reference"} }
	bad := make(map[string]struct{}, len(consistency.Mismatches))
	for _, id := range consistency.Mismatches { bad[id] = struct{}{} }
	plan := MeshReconcilePlan{Reference: candidates[0].ID, Reason: "provider snapshot drift"}
	for _, node := range candidates[1:] {
		if _, ok := bad[string(node.ID)]; ok { plan.Targets = append(plan.Targets, node.ID) }
	}
	return plan
}
