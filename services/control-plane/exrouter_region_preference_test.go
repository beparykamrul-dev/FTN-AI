package main

import "testing"

func TestExRouteScoreRegionPreference(t *testing.T) {
 n:=Node{Healthy:true,LatencyMs:10,PacketLoss:0,Region:"BD"}
 base:=exRouteScore(n,ExRouterRequest{})
 preferred:=exRouteScore(n,ExRouterRequest{Region:"bd"})
 if preferred<=base{t.Fatalf("region preference not applied: base=%v preferred=%v",base,preferred)}
}
