package main

import "testing"

func TestExRouteScoreFormula(t *testing.T) {
 n:=Node{Healthy:true,LatencyMs:20,PacketLoss:2}
 got:=exRouteScore(n,ExRouterRequest{})
 want:=72.0
 if got!=want{t.Fatalf("got %v want %v",got,want)}
}
