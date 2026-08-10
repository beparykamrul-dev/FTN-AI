package mesh

import "sort"

type Path struct {
	NextHop string `json:"next_hop"`
	Metric uint32 `json:"metric"`
	Hops int `json:"hops"`
}

// SelectPaths returns up to maxPaths healthy paths, ordered by metric. Equal
// metrics are retained so a dataplane can later implement ECMP safely.
func SelectPaths(paths []Path, maxPaths int) []Path {
	if maxPaths <= 0 || len(paths) == 0 { return nil }
	out := append([]Path(nil), paths...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Metric == out[j].Metric { return out[i].Hops < out[j].Hops }
		return out[i].Metric < out[j].Metric
	})
	if len(out) > maxPaths { out = out[:maxPaths] }
	return out
}

// RemoveUnhealthy filters paths that are no longer eligible for forwarding.
func RemoveUnhealthy(paths []Path, healthyNextHops map[string]bool) []Path {
	out := make([]Path, 0, len(paths))
	for _, p := range paths {
		if healthyNextHops[p.NextHop] { out = append(out, p) }
	}
	return out
}
