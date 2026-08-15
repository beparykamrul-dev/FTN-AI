package dynamic

import "sort"

// Endpoint is a currently advertised FTNDNS endpoint.
type Endpoint struct {
	Name      string
	Address   string
	NodeID    string
	Healthy   bool
	LatencyMS uint32
	Priority  uint16
}

// SelectHealthy returns deterministic, health-aware endpoints. Lower
// priority and latency win; unhealthy endpoints are excluded.
func SelectHealthy(endpoints []Endpoint) []Endpoint {
	out := make([]Endpoint, 0, len(endpoints))
	for _, e := range endpoints {
		if e.Name != "" && e.Address != "" && e.NodeID != "" && e.Healthy {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority != out[j].Priority {
			return out[i].Priority < out[j].Priority
		}
		if out[i].LatencyMS != out[j].LatencyMS {
			return out[i].LatencyMS < out[j].LatencyMS
		}
		return out[i].NodeID < out[j].NodeID
	})
	return out
}
