package main

import "testing"

func TestSelectFTNECMPPathsDeterministicAndBFDReady(t *testing.T) {
	candidates := []FTNRouteCandidate{
		{PathID:"path-b", Route:FTNRoute{Prefix:"203.0.113.9/24", NextHop:"192.0.2.2", Protocol:"bgp", Active:true}, Healthy:true, BFDState:FTNBFDUp, Score:90},
		{PathID:"path-down", Route:FTNRoute{Prefix:"203.0.113.9/24", NextHop:"192.0.2.3", Protocol:"bgp", Active:true}, Healthy:true, BFDState:FTNBFDDown, Score:99},
		{PathID:"path-a", Route:FTNRoute{Prefix:"203.0.113.9/24", NextHop:"192.0.2.1", Protocol:"bgp", Active:true}, Healthy:true, BFDState:FTNBFDUp, Score:90},
	}
	got, err := SelectFTNECMPPaths(candidates, 2)
	if err != nil { t.Fatal(err) }
	if len(got.Paths) != 2 { t.Fatalf("paths=%d", len(got.Paths)) }
	if got.Paths[0].PathID != "path-a" || got.Paths[1].PathID != "path-b" { t.Fatalf("unexpected order: %+v", got.Paths) }
}

func TestSelectFTNECMPPathsRejectsNoEligiblePath(t *testing.T) {
	_, err := SelectFTNECMPPaths([]FTNRouteCandidate{{PathID:"p", Route:FTNRoute{Prefix:"198.51.100.0/24", NextHop:"192.0.2.1", Protocol:"bgp"}, Healthy:false, BFDState:FTNBFDDown}}, 2)
	if err == nil { t.Fatal("expected no eligible path error") }
}

func TestBuildFTNECMPReconcilePlanDetectsDrift(t *testing.T) {
	current := FTNECMPSelection{Prefix:"203.0.113.0/24", VRF:"default", Paths:[]FTNRouteCandidate{{PathID:"a", Route:FTNRoute{Prefix:"203.0.113.0/24", NextHop:"192.0.2.1", Protocol:"bgp"}}}}
	desired := FTNECMPSelection{Prefix:"203.0.113.0/24", VRF:"default", Paths:[]FTNRouteCandidate{{PathID:"b", Route:FTNRoute{Prefix:"203.0.113.0/24", NextHop:"192.0.2.2", Protocol:"bgp"}}}}
	plan := BuildFTNECMPReconcilePlan(current, desired)
	if !plan.Changed || plan.Reason != "ecmp_state_drift_requires_approval" { t.Fatalf("unexpected plan: %+v", plan) }
}

func TestBuildFTNECMPReconcilePlanNoOp(t *testing.T) {
	p := FTNRouteCandidate{PathID:"a", Route:FTNRoute{Prefix:"203.0.113.0/24", NextHop:"192.0.2.1", Protocol:"bgp"}}
	current := FTNECMPSelection{Prefix:"203.0.113.0/24", VRF:"default", Paths:[]FTNRouteCandidate{p}}
	plan := BuildFTNECMPReconcilePlan(current, current)
	if plan.Changed || plan.Reason != "ecmp_state_in_sync" { t.Fatalf("unexpected plan: %+v", plan) }
}
