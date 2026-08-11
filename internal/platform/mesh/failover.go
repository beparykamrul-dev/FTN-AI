package mesh

import (
	"sort"
	"time"
)

type Path struct { NextHop string `json:"next_hop"`; Metric uint32 `json:"metric"`; Hops int `json:"hops"` }

func SelectPaths(paths []Path, maxPaths int) []Path {
	if maxPaths <= 0 || len(paths) == 0 { return nil }
	out := append([]Path(nil), paths...)
	sort.SliceStable(out, func(i, j int) bool { if out[i].Metric == out[j].Metric { return out[i].Hops < out[j].Hops }; return out[i].Metric < out[j].Metric })
	if len(out) > maxPaths { out = out[:maxPaths] }
	return out
}

func RemoveUnhealthy(paths []Path, healthyNextHops map[string]bool) []Path {
	out := make([]Path, 0, len(paths)); for _, p := range paths { if healthyNextHops[p.NextHop] { out = append(out, p) } }; return out
}

type FailoverDecision struct { CurrentLink string `json:"current_link"`; CandidateLink string `json:"candidate_link"`; Reason string `json:"reason"`; Score uint32 `json:"score"`; GeneratedAt time.Time `json:"generated_at"` }

// ChooseFailover proposes an alternative only; policy, approval, execution and verification remain separate steps.
func ChooseFailover(current string, links []Link, now time.Time, maxAge time.Duration) (FailoverDecision, bool) {
	var best Link; found := false
	for _, l := range links { if l.ID == current || l.State != LinkUp || l.UpdatedAt.IsZero() || now.Sub(l.UpdatedAt) > maxAge { continue }; if !found || l.Metric < best.Metric { best, found = l, true } }
	if !found { return FailoverDecision{}, false }
	return FailoverDecision{CurrentLink: current, CandidateLink: best.ID, Reason: "healthy alternative with lower path metric", Score: best.Metric, GeneratedAt: now.UTC()}, true
}
