package fiber

import "sort"

// BuildTopologyEdges derives relationships from discovered parent IDs. It
// creates a deterministic graph view without changing device configuration.
func BuildTopologyEdges(snapshot TopologySnapshot) []TopologyEdge {
	out := make([]TopologyEdge, 0)
	for _, e := range snapshot.Entities {
		if e.ExternalID == "" || e.ParentID == "" { continue }
		out = append(out, TopologyEdge{From:e.ParentID, To:e.ExternalID, Kind:e.Kind, Source:e.Source})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].From == out[j].From { return out[i].To < out[j].To }
		return out[i].From < out[j].From
	})
	return out
}
