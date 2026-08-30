package main

import "time"

type UnifiedRouteCandidate struct {
    Node Node `json:"node"`
    Decision AdvancedRouteDecision `json:"decision"`
    Fabric RouteFabricState `json:"fabric"`
}

func buildUnifiedCandidates(nodes []Node, req PlacementRequest, policy AdvancedRoutePolicy, now time.Time) []UnifiedRouteCandidate {
    out := make([]UnifiedRouteCandidate, 0, len(nodes))
    for _, n := range nodes {
        if !n.Healthy || !nodeHasService(n, req.ServiceID) { continue }
        fabric := nodeFabricState(n)
        score, eligible, reasons := routeFabricScore(fabric, policy)
        if eligible { score += routePreferenceMatch(n.Provider, n.Region, req.Provider, req.Region) }
        out = append(out, UnifiedRouteCandidate{Node:n, Fabric:fabric, Decision:AdvancedRouteDecision{PathID:n.ID, Score:score, Eligible:eligible, Reasons:reasons}})
    }
    return out
}

func chooseUnifiedRoute(service string, candidates []UnifiedRouteCandidate, policy AdvancedRoutePolicy, now time.Time) (UnifiedRouteCandidate, bool) {
    decisions := make([]AdvancedRouteDecision, 0, len(candidates))
    for _, c := range candidates { decisions = append(decisions, c.Decision) }
    selected, ok := selectAdvancedRoute(service, decisions, policy, now)
    if !ok { return UnifiedRouteCandidate{}, false }
    for _, c := range candidates { if c.Node.ID == selected.PathID { return c, true } }
    return UnifiedRouteCandidate{}, false
}
