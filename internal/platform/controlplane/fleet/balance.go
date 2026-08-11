package fleet

import "sort"

type Node struct {
	ID string `json:"id"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkPercent float64 `json:"network_percent"`
	ServiceHealth float64 `json:"service_health"`
	DNSLoad float64 `json:"dns_load"`
}

type Score struct {
	NodeID string `json:"node_id"`
	Value float64 `json:"value"`
}

// ScoreNodes ranks nodes for capacity-aware placement. It is advisory only:
// execution remains behind the FTN approval/policy layer.
func ScoreNodes(nodes []Node) []Score {
	out := make([]Score, 0, len(nodes))
	for _, n := range nodes {
		load := n.CPUPercent*0.25 + n.MemoryPercent*0.25 + n.NetworkPercent*0.20 + n.DNSLoad*0.15 + (100-n.ServiceHealth)*0.15
		out = append(out, Score{NodeID: n.ID, Value: load})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Value < out[j].Value })
	return out
}
