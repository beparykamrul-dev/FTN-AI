package main

import (
    "testing"
    "time"
)

func TestUnifiedRouteSelection(t *testing.T) {
    now := time.Now().UTC()
    nodes := []Node{
        {ID:"primary", Provider:"p1", Region:"bd", Services:[]string{"edge"}, Healthy:true, BGPUp:true, BFDUp:true, RPKIValid:true, AnycastReady:true, CapacityMbps:10000, UtilizationPercent:30, LastSeen:now},
        {ID:"backup", Provider:"p2", Region:"sg", Services:[]string{"edge"}, Healthy:true, BGPUp:true, BFDUp:true, RPKIValid:true, CapacityMbps:5000, UtilizationPercent:70, LastSeen:now},
    }
    policy := AdvancedRoutePolicy{RequireBGP:true, RequireBFD:true, RequireRPKI:true, PreferAnycast:true, MinScore:1}
    c := buildUnifiedCandidates(nodes, PlacementRequest{ServiceID:"edge", Region:"bd", Provider:"p1"}, policy, now)
    if len(c) != 2 { t.Fatalf("want 2 candidates, got %d", len(c)) }
    best, ok := chooseUnifiedRoute("edge", c, policy, now)
    if !ok || best.Node.ID != "primary" { t.Fatalf("want primary route, got %+v ok=%v", best.Node.ID, ok) }
}

func TestUnifiedRouteRejectsUnsafeFabric(t *testing.T) {
    now := time.Now().UTC()
    nodes := []Node{{ID:"bad", Provider:"p", Services:[]string{"edge"}, Healthy:true, BGPUp:false, BFDUp:true, RPKIValid:true, LastSeen:now}}
    policy := AdvancedRoutePolicy{RequireBGP:true, RequireBFD:true, RequireRPKI:true, MinScore:1}
    c := buildUnifiedCandidates(nodes, PlacementRequest{ServiceID:"edge"}, policy, now)
    if len(c) != 1 || c[0].Decision.Eligible { t.Fatal("unsafe BGP state must be ineligible") }
}
