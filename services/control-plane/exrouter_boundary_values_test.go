package main

import "testing"

func TestExRouteScoreAcceptsBoundaryPacketLoss(t *testing.T) {
 for _,loss:=range []float64{0,100}{n:=Node{Healthy:true,LatencyMs:0,PacketLoss:loss};if got:=exRouteScore(n,ExRouterRequest{});got<0{t.Fatalf("boundary loss %v rejected: %v",loss,got)}}
}
