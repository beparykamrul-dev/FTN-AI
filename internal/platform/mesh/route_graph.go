package mesh

import (
	"math"
	"sort"
)

type RouteHop struct { PeerID string `json:"peer_id"`; Cost float64 `json:"cost"` }
type RouteGraph struct { edges map[string][]RouteHop }
func NewRouteGraph() *RouteGraph { return &RouteGraph{edges: make(map[string][]RouteHop)} }
func (g *RouteGraph) Connect(from string, hop RouteHop) { if g == nil || from == "" || hop.PeerID == "" || from == hop.PeerID || math.IsNaN(hop.Cost) || math.IsInf(hop.Cost, 0) || hop.Cost < 0 { return }; g.edges[from] = append(g.edges[from], hop) }
func (g *RouteGraph) BestPath(source, destination string) ([]string, float64, bool) {
	if g == nil || source == "" || destination == "" { return nil, 0, false }
	dist := map[string]float64{source: 0}; prev := make(map[string]string); visited := make(map[string]bool)
	for { current := ""; best := 0.0; for node, d := range dist { if !visited[node] && (current == "" || d < best || (d == best && node < current)) { current, best = node, d } }; if current == "" { break }; if current == destination { break }; visited[current] = true; for _, edge := range g.edges[current] { if edge.Cost < 0 || math.IsNaN(edge.Cost) || math.IsInf(edge.Cost,0) { continue }; candidate := best + edge.Cost; old, ok := dist[edge.PeerID]; if !ok || candidate < old || (candidate == old && current < prev[edge.PeerID]) { dist[edge.PeerID] = candidate; prev[edge.PeerID] = current } } }
	cost, ok := dist[destination]; if !ok { return nil, 0, false }; path := []string{destination}; for path[len(path)-1] != source { p, ok := prev[path[len(path)-1]]; if !ok { return nil, 0, false }; path = append(path, p) }; for i,j:=0,len(path)-1;i<j;i,j=i+1,j-1 { path[i],path[j]=path[j],path[i] }; return path,cost,true
}
func (g *RouteGraph) Peers(node string) []RouteHop { if g == nil { return []RouteHop{} }; out:=append([]RouteHop(nil),g.edges[node]...); sort.SliceStable(out,func(i,j int)bool{if out[i].Cost!=out[j].Cost{return out[i].Cost<out[j].Cost};return out[i].PeerID<out[j].PeerID}); return out }
