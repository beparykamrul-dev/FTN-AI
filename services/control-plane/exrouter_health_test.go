package main

import "testing"

func TestExRouteScoreRejectsUnhealthyNode(t *testing.T) {
 n:=Node{Healthy:false,LatencyMs:1,PacketLoss:0}
 if got:=exRouteScore(n,ExRouterRequest{});got!=-1{t.Fatalf("got %v",got)}
}
