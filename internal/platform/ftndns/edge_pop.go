package ftndns

import "sort"

type EdgePOP struct {
	ID string `json:"id"`
	Region string `json:"region"`
	Healthy bool `json:"healthy"`
	LatencyMS int64 `json:"latency_ms"`
	Capacity int `json:"capacity"`
}

func RankEdgePOPs(pops []EdgePOP) []EdgePOP {
	ordered := append([]EdgePOP(nil), pops...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Healthy != ordered[j].Healthy { return ordered[i].Healthy }
		if ordered[i].LatencyMS != ordered[j].LatencyMS { return ordered[i].LatencyMS < ordered[j].LatencyMS }
		if ordered[i].Capacity != ordered[j].Capacity { return ordered[i].Capacity > ordered[j].Capacity }
		if ordered[i].Region != ordered[j].Region { return ordered[i].Region < ordered[j].Region }
		return ordered[i].ID < ordered[j].ID
	})
	return ordered
}
