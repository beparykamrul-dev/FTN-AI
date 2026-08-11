package ftndns

import "sort"

// RankMeshNodes returns deterministic operational priority for DNS mesh nodes.
// Healthy nodes are preferred, then lower latency, then stable node ID order.
func RankMeshNodes(nodes []MeshNode) []MeshNode {
	ordered := append([]MeshNode(nil), nodes...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Healthy != ordered[j].Healthy { return ordered[i].Healthy }
		if ordered[i].LatencyMS != ordered[j].LatencyMS { return ordered[i].LatencyMS < ordered[j].LatencyMS }
		return ordered[i].ID < ordered[j].ID
	})
	return ordered
}

func HealthyMeshNodes(nodes []MeshNode) []MeshNode {
	ordered := RankMeshNodes(nodes)
	result := make([]MeshNode, 0, len(ordered))
	for _, node := range ordered { if node.Healthy { result = append(result, node) } }
	return result
}
