package mesh

import "sort"

type RouteHop struct {
	PeerID string `json:"peer_id"`
	Cost float64 `json:"cost"`
}

type RouteGraph struct {
	edges map[string][]RouteHop
}

func NewRouteGraph() *RouteGraph { return &RouteGraph{edges: make(map[string][]RouteHop)} }

func (g *RouteGraph) Connect(from string, hop RouteHop) {
	g.edges[from] = append(g.edges[from], hop)
}

// BestPath returns a minimum-cost path using Dijkstra's algorithm. Invalid
// negative edge costs are ignored because routing metrics are non-negative.
func (g *RouteGraph) BestPath(source, destination string) ([]string, float64, bool) {
	dist := map[string]float64{source: 0}
	prev := make(map[string]string)
	visited := make(map[string]bool)

	for {
		current := ""
		best := 0.0
		for node, d := range dist {
			if !visited[node] && (current == "" || d < best) { current, best = node, d }
		}
		if current == "" { break }
		if current == destination { break }
		visited[current] = true
		for _, edge := range g.edges[current] {
			if edge.Cost < 0 { continue }
			candidate := best + edge.Cost
			old, ok := dist[edge.PeerID]
			if !ok || candidate < old {
				dist[edge.PeerID] = candidate
				prev[edge.PeerID] = current
			}
		}
	}

	cost, ok := dist[destination]
	if !ok { return nil, 0, false }
	path := []string{destination}
	for path[len(path)-1] != source {
		p, ok := prev[path[len(path)-1]]
		if !ok { return nil, 0, false }
		path = append(path, p)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 { path[i], path[j] = path[j], path[i] }
	return path, cost, true
}

func (g *RouteGraph) Peers(node string) []RouteHop {
	out := append([]RouteHop(nil), g.edges[node]...)
	sort.Slice(out, func(i, j int) bool { return out[i].Cost < out[j].Cost })
	return out
}
