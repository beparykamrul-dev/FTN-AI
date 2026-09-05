package fiber

import "sort"

type BoundingBox struct {
	MinLat float64 `json:"min_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLat float64 `json:"max_lat"`
	MaxLon float64 `json:"max_lon"`
}

func (b BoundingBox) Contains(lat, lon float64) bool {
	return b.MinLat <= b.MaxLat && b.MinLon <= b.MaxLon && lat >= b.MinLat && lat <= b.MaxLat && lon >= b.MinLon && lon <= b.MaxLon
}

// NodesInBounds returns a stable GIS query result without modifying topology.
func (t *Topology) NodesInBounds(box BoundingBox) []TopologyNode {
	if t == nil || box.MinLat > box.MaxLat || box.MinLon > box.MaxLon {
		return []TopologyNode{}
	}
	nodes, _ := t.Snapshot()
	out := make([]TopologyNode, 0)
	for _, n := range nodes {
		if box.Contains(n.Latitude, n.Longitude) {
			out = append(out, n)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
