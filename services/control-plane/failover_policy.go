package main

import (
    "sort"
    "time"
)

type FailoverPath struct {
    PathID string `json:"path_id"`
    Provider string `json:"provider"`
    Region string `json:"region"`
    Priority int `json:"priority"`
    Healthy bool `json:"healthy"`
    BGPUp bool `json:"bgp_up"`
    BFDUp bool `json:"bfd_up"`
    AnycastReady bool `json:"anycast_ready"`
    LastSeen time.Time `json:"last_seen"`
}

type FailoverPolicy struct {
    MaxStale time.Duration `json:"max_stale"`
    RequireBGP bool `json:"require_bgp"`
    RequireBFD bool `json:"require_bfd"`
    PreferDistinctProvider bool `json:"prefer_distinct_provider"`
}

func rankFailoverPaths(paths []FailoverPath, policy FailoverPolicy, now time.Time) []FailoverPath {
    if policy.MaxStale <= 0 { policy.MaxStale = 90 * time.Second }
    eligible := make([]FailoverPath, 0, len(paths))
    for _, p := range paths {
        if !p.Healthy || p.LastSeen.IsZero() || now.Sub(p.LastSeen) > policy.MaxStale { continue }
        if policy.RequireBGP && !p.BGPUp { continue }
        if policy.RequireBFD && !p.BFDUp { continue }
        eligible = append(eligible, p)
    }
    sort.SliceStable(eligible, func(i, j int) bool {
        ai, aj := eligible[i], eligible[j]
        if ai.Priority != aj.Priority { return ai.Priority < aj.Priority }
        if ai.AnycastReady != aj.AnycastReady { return ai.AnycastReady }
        if ai.BFDUp != aj.BFDUp { return ai.BFDUp }
        if ai.BGPUp != aj.BGPUp { return ai.BGPUp }
        return ai.PathID < aj.PathID
    })
    return eligible
}

func chooseFailoverPath(paths []FailoverPath, policy FailoverPolicy, now time.Time) (FailoverPath, bool) {
    ranked := rankFailoverPaths(paths, policy, now)
    if len(ranked) == 0 { return FailoverPath{}, false }
    return ranked[0], true
}
