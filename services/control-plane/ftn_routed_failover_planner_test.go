package main

import "testing"

func TestPlanFTNRoutedFailoverSelectsHealthyCore(t *testing.T) {
	nodes := []FTNCoreNode{
		{ID: "core-b", Healthy: true, BGPReady: true, BFDState: FTNBFDUp},
		{ID: "core-a", Healthy: false, BGPReady: false, BFDState: FTNBFDDown},
	}
	plan := PlanFTNRoutedFailover(nodes, "core-a")
	if !plan.Allowed || !plan.RequiresApproval || plan.TargetNode != "core-b" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	if plan.Risk != "high" || plan.DecisionHash == "" {
		t.Fatalf("unexpected risk/hash: %+v", plan)
	}
}

func TestPlanFTNRoutedFailoverNoHealthyCore(t *testing.T) {
	plan := PlanFTNRoutedFailover([]FTNCoreNode{{ID: "core-a", Healthy: false, BGPReady: false, BFDState: FTNBFDDown}}, "core-a")
	if plan.Allowed || plan.TargetNode != "" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}

func TestPlanFTNRoutedFailoverDeterministic(t *testing.T) {
	nodes := []FTNCoreNode{
		{ID: "core-b", Healthy: true, BGPReady: true, BFDState: FTNBFDUp},
		{ID: "core-a", Healthy: false, BGPReady: false, BFDState: FTNBFDDown},
	}
	a := PlanFTNRoutedFailover(nodes, "core-a")
	b := PlanFTNRoutedFailover(nodes, "core-a")
	if a.DecisionHash != b.DecisionHash {
		t.Fatalf("hash changed: %s != %s", a.DecisionHash, b.DecisionHash)
	}
}
