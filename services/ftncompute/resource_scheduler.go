package ftncompute

import "sort"

// Node describes schedulable FTN compute capacity.
type Node struct {
	ID          string
	Healthy     bool
	CPUFree     uint64
	MemoryFree  uint64
	StorageFree uint64
	LatencyMS   uint32
	Load        uint8
}

// Workload describes minimum resources required by a service/VM/CT workload.
type Workload struct {
	ID          string
	CPU         uint64
	Memory      uint64
	Storage     uint64
	MaxLatency  uint32
}

// SelectNode chooses the healthiest node with enough capacity using a
// deterministic score: lower load/latency wins, then higher free resources.
func SelectNode(nodes []Node, w Workload) (Node, bool) {
	candidates := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if !n.Healthy || n.CPUFree < w.CPU || n.MemoryFree < w.Memory || n.StorageFree < w.Storage {
			continue
		}
		if w.MaxLatency > 0 && n.LatencyMS > w.MaxLatency {
			continue
		}
		candidates = append(candidates, n)
	}
	if len(candidates) == 0 {
		return Node{}, false
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Load != candidates[j].Load {
			return candidates[i].Load < candidates[j].Load
		}
		if candidates[i].LatencyMS != candidates[j].LatencyMS {
			return candidates[i].LatencyMS < candidates[j].LatencyMS
		}
		if candidates[i].CPUFree != candidates[j].CPUFree {
			return candidates[i].CPUFree > candidates[j].CPUFree
		}
		return candidates[i].ID < candidates[j].ID
	})
	return candidates[0], true
}
