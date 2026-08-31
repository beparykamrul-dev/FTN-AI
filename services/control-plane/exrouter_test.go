package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExRouteScorePrefersHealthyServicePath(t *testing.T) {
	req := ExRouterRequest{ServiceID: "internet"}
	fast := Node{ID:"fast", Provider:"p1", Region:"r1", Services:[]string{"internet"}, LatencyMs:8, PacketLoss:0.05, JitterMs:1, Retransmissions:0, NetMbps:1000, Healthy:true}
	slow := Node{ID:"slow", Provider:"p2", Region:"r2", Services:[]string{"internet"}, LatencyMs:30, PacketLoss:0.2, JitterMs:5, Retransmissions:2, NetMbps:100, Healthy:true}
	if exRouteScore(fast, req) <= exRouteScore(slow, req) { t.Fatal("expected fast/low-loss path to score higher") }
}

func TestExRouteRejectsNodeWithoutService(t *testing.T) {
	req := ExRouterRequest{ServiceID:"internet"}
	n := Node{ID:"other", Provider:"p1", Services:[]string{"dns"}, Healthy:true}
	if exRouteScore(n, req) >= 0 { t.Fatal("expected service-ineligible node to be rejected") }
}

func TestNetworkHealthSortsHealthyPathsFirst(t *testing.T) {
	old := nodes
	nodes=[]Node{{ID:"stale",Provider:"p1",Services:[]string{"internet"},LatencyMs:1,Healthy:true,LastSeen:time.Now().UTC().Add(-2*nodeHeartbeatTTL)},{ID:"good",Provider:"p2",Services:[]string{"internet"},LatencyMs:10,PacketLoss:.1,Healthy:true,LastSeen:time.Now().UTC()}}
	defer func(){nodes=old}()
	r:=httptest.NewRequest("GET","/api/v1/network/health",nil); w:=httptest.NewRecorder(); (&App{}).networkHealth(w,r)
	if w.Code!=200 {t.Fatalf("status=%d",w.Code)}
	if !strings.Contains(w.Body.String(),"good") {t.Fatal("expected healthy path in response")}
}
