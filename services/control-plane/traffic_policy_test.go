package main

import (
	"testing"
	"time"
)

func TestSelectTrafficPathRejectsUnhealthy(t *testing.T) {
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	decision, ok := SelectTrafficPath([]TrafficPathObservation{{PathID: "p1", ServiceID: "pubg", Class: TrafficGaming, Healthy: false, ObservedAt: now}}, service, now)
	if ok || decision.PathID != "" {
		t.Fatal("expected no usable path")
	}
}

func TestSelectTrafficPathDoesNotCrossMatchSameClass(t *testing.T) {
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	observations := []TrafficPathObservation{
		{PathID: "freefire-path", ServiceID: "freefire", Class: TrafficGaming, LatencyMs: 1, Healthy: true, ObservedAt: now},
		{PathID: "generic-path", ServiceID: "realtime-generic", Class: TrafficRealtime, LatencyMs: 2, Healthy: true, ObservedAt: now},
	}
	if decision, ok := SelectTrafficPath(observations, service, now); ok || decision.PathID != "" {
		t.Fatalf("unexpected cross-service match: %+v", decision)
	}
}

func TestSelectTrafficPathUsesExplicitGenericFallback(t *testing.T) {
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	decision, ok := SelectTrafficPath([]TrafficPathObservation{{PathID: "generic-gaming", ServiceID: "gaming-generic", Class: TrafficGaming, LatencyMs: 5, Healthy: true, ObservedAt: now}}, service, now)
	if ok || decision.PathID != "" {
		t.Fatalf("unexpected unknown generic fallback: %+v", decision)
	}
	decision, ok = SelectTrafficPath([]TrafficPathObservation{{PathID: "generic-gaming", ServiceID: "realtime-generic", Class: TrafficGaming, LatencyMs: 5, Healthy: true, ObservedAt: now}}, service, now)
	if ok || decision.PathID != "" {
		t.Fatalf("unexpected class-mismatched generic fallback: %+v", decision)
	}
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
	if !ok || first.PathID != "p1" {
		t.Fatalf("initial path=%q", first.PathID)
	}
	second, ok := c.Decide(paths, service, base.Add(1*time.Second))
	if !ok || second.PathID != "p1" {
		t.Fatalf("hold-down path=%q", second.PathID)
	}
	third, ok := c.Decide(paths, service, base.Add(6*time.Second))
	if !ok || third.PathID != "p2" {
		t.Fatalf("expected switch to p2, got %q", third.PathID)
	}
}

func TestRouterOSTrafficQoSPlanIsApprovalGated(t *testing.T) {
	device := NetworkDevice{ID: "r1", Kind: "router", Address: "https://router", Healthy: true}
	plan, err := BuildRouterOSTrafficQoSPlan(device, []TrafficDecision{{ServiceID: "whatsapp", Class: TrafficRealtime, DSCP: 46, Priority: 90, PathID: "p1"}})
	if err != nil {
		t.Fatal(err)
	}
	if !plan.RequiresApproval || len(plan.Rules) != 1 {
		t.Fatalf("invalid plan: %+v", plan)
	}
	intent := RouterOSTrafficQoSAction(device)
	decision := EvaluateNetworkExecution(intent)
	if !decision.RequiresApproval || decision.Allowed {
		t.Fatalf("expected approval gate: %+v", decision)
	}
}
