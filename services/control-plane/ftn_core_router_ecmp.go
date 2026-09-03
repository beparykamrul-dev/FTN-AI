package main

import (
	"errors"
	"sort"
	"strings"
)

type FTNRouteCandidate struct {
	PathID string `json:"path_id"`
	Route FTNRoute `json:"route"`
	Healthy bool `json:"healthy"`
	BFDState FTNBFDState `json:"bfd_state"`
	Score float64 `json:"score"`
}

type FTNECMPSelection struct {
	Prefix string `json:"prefix"`
	VRF string `json:"vrf"`
	Paths []FTNRouteCandidate `json:"paths"`
	Changed bool `json:"changed"`
	Reason string `json:"reason"`
}

func normalizeFTNRouteCandidate(c FTNRouteCandidate) (FTNRouteCandidate, error) {
	c.PathID = strings.TrimSpace(c.PathID)
	if c.PathID == "" { return FTNRouteCandidate{}, errors.New("path_id_required") }
	r, err := NormalizeFTNRoute(c.Route)
	if err != nil { return FTNRouteCandidate{}, err }
	if r.NextHop == "" { return FTNRouteCandidate{}, errors.New("next-hop_required") }
	c.Route = r
	if c.BFDState == "" { c.BFDState = FTNBFDUnknown }
	return c, nil
}

func SelectFTNECMPPaths(candidates []FTNRouteCandidate, maxPaths int) (FTNECMPSelection, error) {
	if len(candidates) == 0 { return FTNECMPSelection{}, errors.New("route_candidates_required") }
	if maxPaths <= 0 { return FTNECMPSelection{}, errors.New("max_paths_required") }
	normalized := make([]FTNRouteCandidate, 0, len(candidates))
	for _, c := range candidates {
		n, err := normalizeFTNRouteCandidate(c)
		if err != nil { return FTNECMPSelection{}, err }
		if !n.Healthy || n.BFDState == FTNBFDDown { continue }
		normalized = append(normalized, n)
	}
	if len(normalized) == 0 { return FTNECMPSelection{Reason:"no_eligible_ecmp_path"}, errors.New("no_eligible_ecmp_path") }
	sort.SliceStable(normalized, func(i, j int) bool {
		if normalized[i].Score != normalized[j].Score { return normalized[i].Score > normalized[j].Score }
		if normalized[i].PathID != normalized[j].PathID { return normalized[i].PathID < normalized[j].PathID }
		return normalized[i].Route.NextHop < normalized[j].Route.NextHop
	})
	if len(normalized) > maxPaths { normalized = normalized[:maxPaths] }
	return FTNECMPSelection{Prefix:normalized[0].Route.Prefix, VRF:normalized[0].Route.VRF, Paths:normalized, Reason:"healthy_bfd_eligible_paths_selected"}, nil
}

func DiffFTNECMPSelection(current, desired FTNECMPSelection) bool {
	if current.Prefix != desired.Prefix || current.VRF != desired.VRF || len(current.Paths) != len(desired.Paths) { return true }
	for i := range current.Paths {
		if current.Paths[i].PathID != desired.Paths[i].PathID || current.Paths[i].Route.NextHop != desired.Paths[i].Route.NextHop || current.Paths[i].Route.Prefix != desired.Paths[i].Route.Prefix { return true }
	}
	return false
}

func BuildFTNECMPReconcilePlan(current, desired FTNECMPSelection) FTNECMPSelection {
	desired.Changed = DiffFTNECMPSelection(current, desired)
	if desired.Changed { desired.Reason = "ecmp_state_drift_requires_approval" } else { desired.Reason = "ecmp_state_in_sync" }
	return desired
}
