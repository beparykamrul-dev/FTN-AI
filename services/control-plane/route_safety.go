package main

import (
    "net/netip"
    "strings"
    "time"
)

type RouteSafetyInput struct {
    Prefix string `json:"prefix"`
    RPKIStatus string `json:"rpki_status"`
    ASPathLength int `json:"as_path_length"`
    MaxASPathLength int `json:"max_as_path_length"`
    LocalPref int `json:"local_pref"`
    MED int `json:"med"`
    BFDUp bool `json:"bfd_up"`
    LastUpdate time.Time `json:"last_update"`
    MaxStale time.Duration `json:"max_stale"`
}

func validateRouteSafety(in RouteSafetyInput) (bool, string) {
    p, err := netip.ParsePrefix(strings.TrimSpace(in.Prefix))
    if err != nil || !p.IsValid() { return false, "invalid_prefix" }
    if in.RPKIStatus != "valid" { return false, "rpki_not_valid" }
    if in.MaxASPathLength > 0 && in.ASPathLength > in.MaxASPathLength { return false, "as_path_too_long" }
    if !in.BFDUp { return false, "bfd_down" }
    if in.MaxStale > 0 && (in.LastUpdate.IsZero() || time.Since(in.LastUpdate) > in.MaxStale) { return false, "route_stale" }
    return true, "ok"
}
