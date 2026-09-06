package main

import "testing"

func TestExRouteScoreCaseInsensitivePreferences(t *testing.T) {
 n:=Node{Healthy:true,LatencyMs:5,PacketLoss:0,Region:"Dhaka",Provider:"Transit-A"}
 if exRouteScore(n,ExRouterRequest{Region:"DHAKA"})<=exRouteScore(n,ExRouterRequest{}){t.Fatal("region matching must be case-insensitive")}
 if exRouteScore(n,ExRouterRequest{Provider:"transit-a"})<=exRouteScore(n,ExRouterRequest{}){t.Fatal("provider matching must be case-insensitive")}
}
