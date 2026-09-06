package main

import "testing"

func TestExRouteScoreProviderPreference(t *testing.T) {
 n:=Node{Healthy:true,LatencyMs:10,PacketLoss:0,Provider:"edge-a"}
 base:=exRouteScore(n,ExRouterRequest{})
 preferred:=exRouteScore(n,ExRouterRequest{Provider:"EDGE-A"})
 if preferred<=base{t.Fatalf("provider preference not applied: base=%v preferred=%v",base,preferred)}
}
