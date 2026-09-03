package main

import "testing"

func TestFTNGoBGPAdapter(t *testing.T) {
	a := NewFTNGoBGPAdapter()
	if a.Name() != "gobgp-v4" || a.Established() { t.Fatal("unexpected initial state") }
	if err := a.UpsertPeer(FTNBGPPeerState{PeerID:"core-b", RemoteAddress:"192.0.2.2", RemoteAS:64512, Established:true}); err != nil { t.Fatal(err) }
	a.SetSessionState(true)
	r := FTNRoute{Prefix:"203.0.113.0/24", NextHop:"192.0.2.2", Protocol:"bgp", VRF:"default", Active:true}
	if err := a.ApplyRoute(r); err != nil { t.Fatal(err) }
	if len(a.Routes()) != 1 || len(a.Peers()) != 1 { t.Fatal("route/peer state missing") }
	if err := a.WithdrawRoute(r); err != nil { t.Fatal(err) }
	if len(a.Routes()) != 0 { t.Fatal("route was not withdrawn") }
}

func TestFTNGoBGPRejectsRouteWithoutSession(t *testing.T) {
	a := NewFTNGoBGPAdapter()
	if err := a.ApplyRoute(FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp"}); err == nil { t.Fatal("expected session error") }
}
