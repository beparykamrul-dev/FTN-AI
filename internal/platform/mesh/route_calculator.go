package mesh

import (
	"math"
	"sort"
	"strings"
)

type Route struct { From string `json:"from"`; To string `json:"to"`; NextHop string `json:"next_hop"`; Metric uint32 `json:"metric"`; Hops int `json:"hops"` }

func CalculateRoutes(source string, links []Link) []Route {
	source = strings.TrimSpace(source)
	if source == "" { return []Route{} }
	adj := make(map[string][]Link)
	for _, l := range links {
		l.From, l.To = strings.TrimSpace(l.From), strings.TrimSpace(l.To)
		if l.State != LinkUp || l.From == "" || l.To == "" || l.From == l.To { continue }
		adj[l.From] = append(adj[l.From], l)
	}
	for node := range adj { sort.SliceStable(adj[node], func(i,j int) bool { if adj[node][i].Metric != adj[node][j].Metric { return adj[node][i].Metric < adj[node][j].Metric }; return adj[node][i].To < adj[node][j].To }) }
	type item struct { node string; cost uint32; hops int; first string }
	best := map[string]item{source: {node: source}}
	queue := []item{{node: source}}
	for len(queue) > 0 {
		cur := queue[0]; queue = queue[1:]
		for _, l := range adj[cur.node] {
			if cur.cost > math.MaxUint32-l.Metric { continue }
			first := l.To; if cur.first != "" { first = cur.first }
			candidate := item{node: l.To, cost: cur.cost + l.Metric, hops: cur.hops + 1, first: first}
			old, ok := best[l.To]
			if !ok || candidate.cost < old.cost || (candidate.cost == old.cost && (candidate.hops < old.hops || (candidate.hops == old.hops && candidate.first < old.first))) { best[l.To] = candidate; queue = append(queue, candidate) }
		}
	}
	out := make([]Route, 0, len(best)-1)
	for node, v := range best { if node != source { out = append(out, Route{From: source, To: node, NextHop: v.first, Metric: v.cost, Hops: v.hops}) } }
	sort.Slice(out, func(i,j int) bool { return out[i].To < out[j].To })
	return out
}
