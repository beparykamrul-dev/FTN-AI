package integration

import "sort"

// RouteCandidate combines DNS health, security and network state for a
// provider/node decision. Higher score is preferred.
type RouteCandidate struct {
	Provider string
	NodeID   string
	Healthy  bool
	Secure   bool
	LatencyMS int64
	Load     float64
}

func RankRoutes(candidates []RouteCandidate) []RouteCandidate {
	out := make([]RouteCandidate, 0, len(candidates))
	for _, c := range candidates {
		if c.Healthy {
			out = append(out, c)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Secure != out[j].Secure { return out[i].Secure }
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		if out[i].Load != out[j].Load { return out[i].Load < out[j].Load }
		if out[i].Provider != out[j].Provider { return out[i].Provider < out[j].Provider }
		return out[i].NodeID < out[j].NodeID
	})
	return out
}
