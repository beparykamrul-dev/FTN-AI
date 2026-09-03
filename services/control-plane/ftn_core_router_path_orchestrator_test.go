package main

import (
	"testing"
	"time"
)

func TestFTNCoreRouterPathOrchestratorUsesHysteresis(t *testing.T) {
	o := NewFTNCoreRouterPathOrchestrator(NewFTNRoutedEngine(nil, nil))
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "pubg", Class: TrafficGaming, Priority: 95, DSCP: 46}
	obs := []TrafficPathObservation{
		{PathID: "a", ServiceID: "pubg", Class: TrafficGaming, LatencyMs: 10, ObservedAt: now, Healthy: true},
		{PathID: "b", ServiceID: "pubg", Class: TrafficGaming, LatencyMs: 20, ObservedAt: now, Healthy: true},
	}
	first, ok := o.SelectPath(obs, service, now)
	if !ok || first.PathID != "a" { t.Fatalf("initial path = %#v, ok=%v", first, ok) }
	obs[1].LatencyMs = 1
	second, ok := o.SelectPath(obs, service, now.Add(time.Second))
	if !ok || second.PathID != "a" { t.Fatalf("hysteresis path = %#v, ok=%v", second, ok) }
	third, ok := o.SelectPath(obs, service, now.Add(7*time.Second))
	if !ok || third.PathID != "b" { t.Fatalf("switched path = %#v, ok=%v", third, ok) }
}

func TestFTNCoreRouterPathOrchestratorFailsOverUnhealthyCurrent(t *testing.T) {
	o := NewFTNCoreRouterPathOrchestrator(NewFTNRoutedEngine(nil, nil))
	now := time.Unix(1000, 0)
	service := TrafficServicePolicy{ID: "whatsapp", Class: TrafficRealtime, Priority: 90, DSCP: 46}
	obs := []TrafficPathObservation{
		{PathID: "a", ServiceID: "whatsapp", Class: TrafficRealtime, LatencyMs: 5, ObservedAt: now, Healthy: true},
		{PathID: "b", ServiceID: "whatsapp", Class: TrafficRealtime, LatencyMs: 15, ObservedAt: now, Healthy: true},
	}
	if _, ok := o.SelectPath(obs, service, now); !ok { t.Fatal("initial selection failed") }
	obs[0].Healthy = false
	got, ok := o.SelectPath(obs, service, now.Add(time.Second))
	if !ok || got.PathID != "b" || !got.Failover { t.Fatalf("failover = %#v, ok=%v", got, ok) }
}

func TestFTNCoreRouterPathOrchestratorRejectsMissingEngine(t *testing.T) {
	o := NewFTNCoreRouterPathOrchestrator(nil)
	_, err := o.BuildECMPPlan(FTNECMPSelection{}, []FTNRouteCandidate{{PathID: "a", Route: FTNRoute{Prefix: "203.0.113.0/24", NextHop: "192.0.2.1", Protocol: "bgp"}, Healthy: true, BFDState: FTNBFDUp}}, 2)
	if err == nil { t.Fatal("expected routed engine error") }
}
