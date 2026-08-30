package main

import (
    "math"
    "strings"
    "sync"
    "time"
)

type RouteFabricState struct {
    BGPUp bool `json:"bgp_up"`
    BFDUp bool `json:"bfd_up"`
    ISISUp bool `json:"isis_up"`
    EVPNReady bool `json:"evpn_ready"`
    AnycastReady bool `json:"anycast_ready"`
    RPKIValid bool `json:"rpki_valid"`
    PrefixCount int `json:"prefix_count"`
    CapacityMbps float64 `json:"capacity_mbps"`
    UtilizationPercent float64 `json:"utilization_percent"`
}

type AdvancedRoutePolicy struct {
    MinScore float64 `json:"min_score"`
    MaxUtilizationPercent float64 `json:"max_utilization_percent"`
    RequireBGP bool `json:"require_bgp"`
    RequireBFD bool `json:"require_bfd"`
    RequireRPKI bool `json:"require_rpki"`
    PreferAnycast bool `json:"prefer_anycast"`
    PreferEVPN bool `json:"prefer_evpn"`
    StickySeconds int `json:"sticky_seconds"`
}

type AdvancedRouteDecision struct {
    PathID string `json:"path_id"`
    Score float64 `json:"score"`
    Eligible bool `json:"eligible"`
    Reasons []string `json:"reasons,omitempty"`
}

var advancedRouteState = struct { sync.Mutex; selected map[string]struct{Path string; At time.Time} }{selected: map[string]struct{Path string; At time.Time}{}}

func normalizePolicy(p AdvancedRoutePolicy) AdvancedRoutePolicy {
    if p.MinScore == 0 { p.MinScore = 1 }
    if p.MaxUtilizationPercent <= 0 || p.MaxUtilizationPercent > 100 { p.MaxUtilizationPercent = 85 }
    if p.StickySeconds <= 0 { p.StickySeconds = 30 }
    return p
}

func routeFabricScore(s RouteFabricState, p AdvancedRoutePolicy) (float64, bool, []string) {
    reasons := []string{}
    if p.RequireBGP && !s.BGPUp { return -1, false, []string{"bgp_down"} }
    if p.RequireBFD && !s.BFDUp { return -1, false, []string{"bfd_down"} }
    if p.RequireRPKI && !s.RPKIValid { return -1, false, []string{"rpki_invalid"} }
    if s.UtilizationPercent > p.MaxUtilizationPercent { return -1, false, []string{"capacity_guard"} }
    score := 0.0
    if s.BGPUp { score += 20 } else { reasons = append(reasons, "bgp_not_confirmed") }
    if s.BFDUp { score += 15 }
    if s.ISISUp { score += 10 }
    if s.EVPNReady { score += 8 }
    if s.AnycastReady { score += 8 }
    if s.RPKIValid { score += 15 }
    score += math.Min(float64(s.PrefixCount)/10000, 5)
    if s.CapacityMbps > 0 { score += math.Min(s.CapacityMbps/10000, 10) }
    score -= s.UtilizationPercent * 0.2
    if p.PreferAnycast && s.AnycastReady { score += 8 }
    if p.PreferEVPN && s.EVPNReady { score += 8 }
    return score, score >= p.MinScore, reasons
}

func selectAdvancedRoute(service string, candidates []AdvancedRouteDecision, policy AdvancedRoutePolicy, now time.Time) (AdvancedRouteDecision, bool) {
    policy = normalizePolicy(policy)
    var best AdvancedRouteDecision
    found := false
    for _, c := range candidates { if c.Eligible && c.Score >= policy.MinScore && (!found || c.Score > best.Score) { best, found = c, true } }
    advancedRouteState.Lock(); defer advancedRouteState.Unlock()
    if old, ok := advancedRouteState.selected[service]; ok && now.Sub(old.At) < time.Duration(policy.StickySeconds)*time.Second {
        for _, c := range candidates { if c.PathID == old.Path && c.Eligible { return c, true } }
    }
    if found { advancedRouteState.selected[service] = struct{Path string; At time.Time}{best.PathID, now} }
    return best, found
}

func routePreferenceMatch(provider, region, wantProvider, wantRegion string) float64 {
    score := 0.0
    if wantProvider != "" && strings.EqualFold(provider, wantProvider) { score += 10 }
    if wantRegion != "" && strings.EqualFold(region, wantRegion) { score += 10 }
    return score
}
