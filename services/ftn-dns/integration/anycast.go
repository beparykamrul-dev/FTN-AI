package integration

import "sort"

// AnycastNode is the normalized state consumed by FTN DNS/SDN routing.
type AnycastNode struct {
	ID        string
	Address   string
	Healthy   bool
	LatencyMS int64
	Load      float64
}

// RankAnycastNodes returns healthy nodes first, preferring lower latency and
// then lower load. It does not modify routing state by itself.
func RankAnycastNodes(nodes []AnycastNode) []AnycastNode {
	out := make([]AnycastNode, 0, len(nodes))
	for _, n := range nodes {
		if n.Healthy {
			out = append(out, n)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].LatencyMS != out[j].LatencyMS {
			return out[i].LatencyMS < out[j].LatencyMS
		}
		if out[i].Load != out[j].Load {
			return out[i].Load < out[j].Load
		}
		return out[i].ID < out[j].ID
	})
	return out
}
