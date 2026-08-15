package ftnmap

// ImpactNode describes an affected upstream service in the dependency graph.
type ImpactNode struct {
	ServiceID string
	DependsOn string
	Depth     uint16
}

// BuildImpact returns a bounded, deterministic dependency impact list.
func BuildImpact(edges map[string][]string, failed string, maxDepth uint16) []ImpactNode {
	if failed == "" { return nil }
	seen := map[string]bool{failed:true}
	queue := []ImpactNode{{ServiceID: failed, Depth: 0}}
	result := make([]ImpactNode, 0)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, upstream := range edges[cur.ServiceID] {
			if seen[upstream] || cur.Depth >= maxDepth { continue }
			seen[upstream] = true
			result = append(result, ImpactNode{ServiceID: upstream, DependsOn: cur.ServiceID, Depth: cur.Depth+1})
			queue = append(queue, ImpactNode{ServiceID: upstream, Depth: cur.Depth+1})
		}
	}
	return result
}
