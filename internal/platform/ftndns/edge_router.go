package ftndns

import "sort"

type EdgeRoute struct { Node EdgeNodeState; Weight int }

// RankEdgeRoutes provides deterministic edge-node preference. It does not
// perform routing or mutate provider state.
func RankEdgeRoutes(routes []EdgeRoute) []EdgeRoute {
	ordered := append([]EdgeRoute(nil), routes...)
	sort.SliceStable(ordered, func(i,j int) bool {
		if ordered[i].Node.Healthy != ordered[j].Node.Healthy { return ordered[i].Node.Healthy }
		if ordered[i].Node.LatencyMS != ordered[j].Node.LatencyMS { return ordered[i].Node.LatencyMS < ordered[j].Node.LatencyMS }
		if ordered[i].Weight != ordered[j].Weight { return ordered[i].Weight > ordered[j].Weight }
		return ordered[i].Node.Provider < ordered[j].Node.Provider
	})
	return ordered
}
