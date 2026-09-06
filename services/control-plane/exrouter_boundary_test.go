package main

import (
	"math"
	"testing"
)

func TestExRouteScoreRejectsInvalidMetrics(t *testing.T) {
	n := Node{ID: "n1", Healthy: true, LatencyMs: math.NaN(), PacketLoss: 0}
	if score := exRouteScore(n, ExRouterRequest{}); score >= 0 {
		t.Fatalf("invalid latency was accepted: %v", score)
	}
}

func TestExRouteScoreRewardsMatchingRegionAndProvider(t *testing.T) {
	n := Node{ID: "n1", Healthy: true, Region: "dhaka", Provider: "ftn", LatencyMs: 5, PacketLoss: 0}
	base := exRouteScore(n, ExRouterRequest{})
	matched := exRouteScore(n, ExRouterRequest{Region: "dhaka", Provider: "ftn"})
	if matched <= base {
		t.Fatalf("expected matching route score to improve: base=%v matched=%v", base, matched)
	}
}
