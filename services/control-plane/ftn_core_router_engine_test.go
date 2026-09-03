package main

import (
    "testing"
)

func TestFTNCoreRouterEngineValidatesAndSnapshotsDeterministically(t *testing.T) {
    node := CoreRouterNode{ID:"core-a", Role:"primary", Protocol:"gobgp-v4", Enabled:true}
    bgp := NewFTNGoBGPAdapter()
    bgp.SetSessionState(true)
    e, err := NewFTNCoreRouterEngine(node, bgp)
    if err != nil { t.Fatal(err) }
    if !e.Ready() { t.Fatal("enabled core with established BGP must be ready") }
    if err := e.UpsertPeer(CoreRouterPeer{ID:"peer-b", LocalNode:"core-a", RemoteAS:64513, RemoteIP:"192.0.2.2", AddressFamily:"ipv4-unicast", Enabled:true}); err != nil { t.Fatal(err) }
    if err := e.UpsertPeer(CoreRouterPeer{ID:"peer-a", LocalNode:"core-a", RemoteAS:64512, RemoteIP:"2001:db8::2", AddressFamily:"ipv6-unicast", Enabled:true}); err != nil { t.Fatal(err) }
    snap := e.Snapshot()
    if len(snap.Peers) != 2 || snap.Peers[0].ID != "peer-a" || snap.Peers[1].ID != "peer-b" { t.Fatalf("peers not stable: %+v", snap.Peers) }
}

func TestFTNCoreRouterEngineRouteLifecycle(t *testing.T) {
    node := CoreRouterNode{ID:"core-a", Role:"primary", Protocol:"gobgp-v4", Enabled:true}
    bgp := NewFTNGoBGPAdapter(); bgp.SetSessionState(true)
    e, err := NewFTNCoreRouterEngine(node, bgp); if err != nil { t.Fatal(err) }
    r := FTNRoute{Prefix:"203.0.113.7/24", NextHop:"192.0.2.1", Protocol:"bgp", VRF:"default", Active:true}
    if err := e.AnnounceRoute(r); err != nil { t.Fatal(err) }
    if len(e.Snapshot().Routes) != 1 || e.Snapshot().Routes[0].Prefix != "203.0.113.0/24" { t.Fatalf("route was not normalized: %+v", e.Snapshot().Routes) }
    if err := e.WithdrawRoute(r); err != nil { t.Fatal(err) }
    if len(e.Snapshot().Routes) != 0 { t.Fatal("withdraw must remove route state") }
}

func TestFTNCoreRouterEngineRejectsUnestablishedBGP(t *testing.T) {
    node := CoreRouterNode{ID:"core-a", Role:"primary", Protocol:"gobgp-v4", Enabled:true}
    bgp := NewFTNGoBGPAdapter()
    e, err := NewFTNCoreRouterEngine(node, bgp); if err != nil { t.Fatal(err) }
    if err := e.AnnounceRoute(FTNRoute{Prefix:"203.0.113.0/24", Protocol:"bgp", Active:true}); err == nil { t.Fatal("route announcement must fail without established BGP") }
}
