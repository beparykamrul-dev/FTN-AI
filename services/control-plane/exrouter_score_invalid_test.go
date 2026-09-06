package main

import "testing"

func TestExRouteScoreRejectsInvalidLatency(t *testing.T) {
 n:=Node{Healthy:true,LatencyMs:-1,PacketLoss:0}
 if got:=exRouteScore(n,ExRouterRequest{});got!=-1{t.Fatalf("got %v",got)}
}
