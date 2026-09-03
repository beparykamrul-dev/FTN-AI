package main

import (
    "testing"
    "time"
)

func TestFTNPathOrchestratorPlanSelectsHealthyPath(t *testing.T) {
    now := time.Unix(1000, 0)
    o := NewFTNPathOrchestrator()
    in := FTNPathOrchestratorInput{
        ServiceID: "whatsapp",
        Observations: []TrafficPathObservation{
            {PathID: "path-b", ServiceID: "whatsapp", Healthy: true, LatencyMs: 20, ObservedAt: now},
            {PathID: "path-a", ServiceID: "whatsapp", Healthy: true, LatencyMs: 5, ObservedAt: now},
        },
        ECMPCandidates: []FTNRouteCandidate{
            {PathID: "path-a", Route: FTNRoute{Prefix: "203.0.113.0/24", NextHop: "192.0.2.1", Protocol: "bgp", VRF: "default"}, Healthy: true, BFDState: FTNBFDUp, Score: 95},
            {PathID: "path-b", Route: FTNRoute{Prefix: "203.0.113.0/24", NextHop: "192.0.2.2", Protocol: "bgp", VRF: "default"}, Healthy: true, BFDState: FTNBFDUp, Score: 80},
        },
        CurrentECMP: FTNECMPSelection{Prefix: "203.0.113.0/24", VRF: "default"},
        MaxECMPPaths: 8,
        CoreNodes: []FTNCoreNode{{ID: "core-a", Healthy: true, BGPReady: true, BFDState: FTNBFDUp}},
        CurrentCoreNode: "core-a",
        Now: now,
    }
    plan, err := o.Plan(in)
    if err != nil { t.Fatal(err) }
    if plan.Traffic.PathID != "path-a" { t.Fatalf("selected path=%q", plan.Traffic.PathID) }
    if !plan.HasTraffic || !plan.HasECMP || !plan.RequiresApproval { t.Fatalf("unexpected plan: %+v", plan) }
}

func TestFTNPathOrchestratorRejectsNoHealthyCore(t *testing.T) {
    o := NewFTNPathOrchestrator()
    _, err := o.Plan(FTNPathOrchestratorInput{
        ServiceID: "telegram",
        Observations: []TrafficPathObservation{{PathID: "p1", ServiceID: "telegram", Healthy: true, ObservedAt: time.Unix(1000, 0)}},
        ECMPCandidates: []FTNRouteCandidate{{PathID: "p1", Route: FTNRoute{Prefix: "203.0.113.0/24", NextHop: "192.0.2.1", Protocol: "bgp", VRF: "default"}, Healthy: true, BFDState: FTNBFDUp}},
        MaxECMPPaths: 8,
        CoreNodes: []FTNCoreNode{{ID: "core-a", Healthy: false, BGPReady: true, BFDState: FTNBFDUp}},
    })
    if err == nil { t.Fatal("expected no healthy core error") }
}

func TestFTNPathOrchestratorRequiresRouteForSelectedTrafficPath(t *testing.T) {
    now := time.Unix(1000, 0)
    o := NewFTNPathOrchestrator()
    _, err := o.Plan(FTNPathOrchestratorInput{
        ServiceID: "imo",
        Observations: []TrafficPathObservation{{PathID: "p1", ServiceID: "imo", Healthy: true, ObservedAt: now}},
        ECMPCandidates: []FTNRouteCandidate{{PathID: "p2", Route: FTNRoute{Prefix: "203.0.113.0/24", NextHop: "192.0.2.2", Protocol: "bgp", VRF: "default"}, Healthy: true, BFDState: FTNBFDUp}},
        MaxECMPPaths: 8,
        CoreNodes: []FTNCoreNode{{ID: "core-a", Healthy: true, BGPReady: true, BFDState: FTNBFDUp}},
        CurrentCoreNode: "core-a",
        Now: now,
    })
    if err == nil { t.Fatal("expected selected path route candidate error") }
}
