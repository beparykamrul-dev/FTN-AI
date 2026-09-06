package main

import "testing"

func TestExRouteScoreRejectsPacketLossOutOfRange(t *testing.T) {
 for _,loss:=range []float64{-0.1,100.1}{n:=Node{Healthy:true,LatencyMs:5,PacketLoss:loss};if got:=exRouteScore(n,ExRouterRequest{});got!=-1{t.Fatalf("loss %v got %v",loss,got)}}
}
