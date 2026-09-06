package main

import "testing"

func TestPlanCoreRouteChangeRequiresApproval(t *testing.T) {
	plan := PlanCoreRouteChange(
		CoreRouterNode{ID: "core-1", Role: "primary", Protocol: "bgp"},
		RouteChangeIntent{Action: "announce", Prefix: "10.0.0.0/24", NextHop: "10.0.0.1"},
	)
	if plan.Allowed || len(plan.ValidationErrors) == 0 {
		t.Fatalf("route change without approval must be blocked: %+v", plan)
	}
}

func TestValidateCoreRouterNodeRejectsInvalidRole(t *testing.T) {
	if err := ValidateCoreRouterNode(CoreRouterNode{ID: "core-1", Role: "edge", Protocol: "bgp"}); err == nil {
		t.Fatal("invalid role must be rejected")
	}
}
