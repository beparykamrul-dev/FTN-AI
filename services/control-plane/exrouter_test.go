package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExRouteScorePrefersLowerLatencyAndLoss(t *testing.T) {
	req := ExRouterRequest{ServiceID: "internet"}
	fast := Node{ID:"fast", Provider:"p1", Region:"r1", LatencyMs:8, PacketLoss:0.05, Healthy:true}
	slow := Node{ID:"slow", Provider:"p2", Region:"r2", LatencyMs:30, PacketLoss:0.2, Healthy:true}
	if exRouteScore(fast, req) <= exRouteScore(slow, req) { t.Fatal("expected fast/low-loss path to score higher") }
}

func TestNetworkHealthSortsHealthyPathsFirst(t *testing.T) {
	old := nodes
	nodes = []Node{
		{ID:"stale", Provider:"p1", LatencyMs:1, PacketLoss:0, Healthy:true, LastSeen:time.Now().UTC().Add(-2 * nodeHeartbeatTTL)},
		{ID:"good", Provider:"p2", LatencyMs:10, PacketLoss:0.1, Healthy:true, LastSeen:time.Now().UTC()},
	}
	defer func(){ nodes = old }()
	r := httptest.NewRequest("GET", "/api/v1/network/health", nil)
	w := httptest.NewRecorder()
	(&App{}).networkHealth(w, r)
	if w.Code != 200 { t.Fatalf("status=%d", w.Code) }
	if !strings.Contains(w.Body.String(), "good") { t.Fatal("expected healthy path in response") }
}
