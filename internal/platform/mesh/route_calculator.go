package mesh

import "sort"

type Route struct {
	From string `json:"from"`
	To string `json:"to"`
	NextHop string `json:"next_hop"`
	Metric uint32 `json:"metric"`
	Hops int `json:"hops"`
}

// CalculateRoutes computes shortest paths over the currently supplied healthy
// links. It is intentionally provider-neutral; a later dataplane adapter can
// translate the resulting routes into Linux, RouterOS, FRR or FTN router state.
func CalculateRoutes(source string, links []Link) []Route {
	adj := make(map[string][]Link)
	for _, l := range links {
		if l.State != LinkUp { continue }
		adj[l.From] = append(adj[l.From], l)
	}
	type item struct { node string; cost uint32; hops int; first string }
	best := map[string]item{source: {node: source}}
	queue := []item{{node: source}}
	for len(queue) > 0 {
		cur := queue[0]; queue = queue[1:]
		for _, l := range adj[cur.node] {
			first := l.To
			if cur.first != "" { first = cur.first }
			candidate := item{node:l.To, cost:cur.cost+l.Metric, hops:cur.hops+1, first:first}
			old, ok := best[l.To]
			if !ok || candidate.cost < old.cost || (candidate.cost == old.cost && candidate.hops < old.hops) {
				best[l.To] = candidate
				queue = append(queue, candidate)
			}
		}
	}
	out := make([]Route, 0, len(best)-1)
	for node, v := range best {
		if node == source { continue }
		out = append(out, Route{From:source, To:node, NextHop:v.first, Metric:v.cost, Hops:v.hops})
	}
	sort.Slice(out, func(i,j int) bool { return out[i].To < out[j].To })
	return out
}
