package ftnwan

import "sort"

// Path describes a candidate FTNWAN path between two nodes.
type Path struct {
	ID        string
	Source    string
	Target    string
	Healthy   bool
	LatencyMS uint32
	LossPPM   uint32
	Capacity  uint64
	Hops      uint16
}

// SelectBest returns healthy paths ordered for deterministic, quality-aware
// selection. Lower loss/latency/hops are preferred; higher capacity breaks ties.
func SelectBest(paths []Path) []Path {
	out := make([]Path, 0, len(paths))
	for _, p := range paths {
		if p.ID != "" && p.Source != "" && p.Target != "" && p.Healthy {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].LossPPM != out[j].LossPPM { return out[i].LossPPM < out[j].LossPPM }
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		if out[i].Hops != out[j].Hops { return out[i].Hops < out[j].Hops }
		if out[i].Capacity != out[j].Capacity { return out[i].Capacity > out[j].Capacity }
		return out[i].ID < out[j].ID
	})
	return out
}
