package main

import (
	"testing"
	"time"
)

func TestSelectTrafficPathRejectsUnhealthy(t *testing.T) {
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	decision, ok := SelectTrafficPath([]TrafficPathObservation{{PathID: "p1", ServiceID: "pubg", Class: TrafficGaming, Healthy: false, ObservedAt: now}}, service, now)
	if ok || decision.PathID != "" { t.Fatal("expected no usable path") }
}

func TestTrafficPathControllerHoldDown(t *testing.T) {
	base := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	paths := []TrafficPathObservation{
		{PathID: "p1", ServiceID: "pubg", Class: TrafficGaming, LatencyMs: 10, Healthy: true, ObservedAt: base},
		{PathID: "p2", ServiceID: "pubg", Class: TrafficGaming, LatencyMs: 1, Healthy: true, ObservedAt: base},
	}
	var c TrafficPathController
	first, ok := c.Decide(paths[:1], service, base)
	if !ok || first.PathID != "p1" { t.Fatalf("initial path=%q", first.PathID) }
	second, ok := c.Decide(paths, service, base.Add(1*time.Second))
	if !ok || second.PathID != "p1" { t.Fatalf("hold-down path=%q", second.PathID) }
	third, ok := c.Decide(paths, service, base.Add(6*time.Second))
	if !ok || third.PathID != "p2" { t.Fatalf("expected switch to p2, got %q", third.PathID) }
}

func TestRouterOSTrafficQoSPlanIsApprovalGated(t *testing.T) {
	device := NetworkDevice{ID: "r1", Kind: "router", Address: "https://router", Healthy: true}
	plan, err := BuildRouterOSTrafficQoSPlan(device, []TrafficDecision{{ServiceID: "whatsapp", Class: TrafficRealtime, DSCP: 46, Priority: 90, PathID: "p1"}})
	if err != nil { t.Fatal(err) }
	if !plan.RequiresApproval || len(plan.Rules) != 1 { t.Fatalf("invalid plan: %+v", plan) }
	intent := RouterOSTrafficQoSAction(device)
	decision := EvaluateNetworkExecution(intent)
	if !decision.RequiresApproval || decision.Allowed { t.Fatalf("expected approval gate: %+v", decision) }
}
