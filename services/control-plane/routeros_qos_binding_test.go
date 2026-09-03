package main

import "testing"

func TestRouterOSPathRegistryRejectsUnknownPath(t *testing.T) {
	_, err := (RouterOSPathRegistry{Bindings: []RouterOSPathBinding{{PathID: "path-a", RouteMark: "ftn-a", QueueClass: "realtime"}}}).Resolve("path-b")
	if err == nil { t.Fatal("expected unknown path rejection") }
}

func TestBuildRouterOSQOSPlanWithRegistry(t *testing.T) {
	device := NetworkDevice{ID: "r1", Kind: "core-router", Healthy: true}
	decisions := []TrafficDecision{{ServiceID: "whatsapp", PathID: "path-a", Class: TrafficClassVoice, DSCP: 46, Priority: 90}}
	registry := RouterOSPathRegistry{Bindings: []RouterOSPathBinding{{PathID: "path-a", RouteMark: "ftn-a", QueueClass: "voice"}}}
	plan, err := BuildRouterOSQOSPlanWithRegistry(device, decisions, registry)
	if err != nil { t.Fatal(err) }
	if len(plan.Rules) != 1 || plan.Rules[0].PathID != "path-a" { t.Fatalf("unexpected plan: %+v", plan) }
}
