package main

import (
    "errors"
    "net/netip"
    "sort"
    "strings"
    "sync"
    "time"
)

// FTNCoreRouterEngine is the native FTN routing state boundary. MikroTik is
// intentionally outside this engine and is handled by its own adapter.
type FTNCoreRouterEngine struct {
    mu sync.RWMutex
    node CoreRouterNode
    bgp *FTNGoBGPAdapter
    peers map[string]CoreRouterPeer
    routes map[string]FTNCoreRouteState
}

type FTNCoreRouteState struct {
    VRF string `json:"vrf,omitempty"`
    Prefix string `json:"prefix"`
    NextHop string `json:"next_hop,omitempty"`
    PeerID string `json:"peer_id,omitempty"`
    Installed bool `json:"installed"`
    UpdatedAt time.Time `json:"updated_at"`
}

type FTNCoreRouterSnapshot struct {
    Node CoreRouterNode `json:"node"`
    Peers []CoreRouterPeer `json:"peers"`
    Routes []FTNCoreRouteState `json:"routes"`
    CapturedAt time.Time `json:"captured_at"`
}

func NewFTNCoreRouterEngine(node CoreRouterNode, bgp *FTNGoBGPAdapter) (*FTNCoreRouterEngine, error) {
    if err := ValidateCoreRouterNode(node); err != nil { return nil, err }
    if bgp == nil { return nil, errors.New("gobgp_adapter_required") }
    return &FTNCoreRouterEngine{node: node, bgp: bgp, peers: make(map[string]CoreRouterPeer), routes: make(map[string]FTNCoreRouteState)}, nil
}

func validateCoreRoute(prefix, nextHop string) error {
    if strings.TrimSpace(prefix) == "" { return errors.New("route_prefix_required") }
    if _, err := netip.ParsePrefix(prefix); err != nil { return errors.New("invalid_route_prefix") }
    if strings.TrimSpace(nextHop) != "" {
        if _, err := netip.ParseAddr(nextHop); err != nil { return errors.New("invalid_route_next_hop") }
    }
    return nil
}

func routeKey(vrf, prefix string) string { return strings.TrimSpace(vrf) + "|" + strings.TrimSpace(prefix) }

func (e *FTNCoreRouterEngine) UpsertPeer(p CoreRouterPeer) error {
    if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.LocalNode) == "" || p.RemoteASN == 0 { return errors.New("invalid_core_peer") }
    if _, err := netip.ParseAddr(strings.TrimSpace(p.RemoteIP)); err != nil { return errors.New("invalid_core_peer_ip") }
    if p.AddressFamily != "ipv4-unicast" && p.AddressFamily != "ipv6-unicast" && p.AddressFamily != "l2vpn-evpn" { return errors.New("unsupported_address_family") }
    if p.LocalNode != e.node.ID { return errors.New("peer_local_node_mismatch") }
    e.mu.Lock(); defer e.mu.Unlock(); e.peers[p.ID] = p; return nil
}

func (e *FTNCoreRouterEngine) AnnounceRoute(r FTNRoute) error {
    normalized, err := NormalizeFTNRoute(r)
    if err != nil { return err }
    if !e.bgp.Established() { return errors.New("core_bgp_session_not_established") }
    if err := e.bgp.ApplyRoute(normalized); err != nil { return err }
    e.mu.Lock(); defer e.mu.Unlock()
    e.routes[routeKey(normalized.VRF, normalized.Prefix)] = FTNCoreRouteState{VRF:normalized.VRF, Prefix:normalized.Prefix, NextHop:normalized.NextHop, Installed:true, UpdatedAt:time.Now().UTC()}
    return nil
}

func (e *FTNCoreRouterEngine) WithdrawRoute(r FTNRoute) error {
    normalized, err := NormalizeFTNRoute(r)
    if err != nil { return err }
    if err := e.bgp.WithdrawRoute(normalized); err != nil { return err }
    e.mu.Lock(); defer e.mu.Unlock(); delete(e.routes, routeKey(normalized.VRF, normalized.Prefix)); return nil
}

func (e *FTNCoreRouterEngine) Snapshot() FTNCoreRouterSnapshot {
    e.mu.RLock(); defer e.mu.RUnlock()
    peers := make([]CoreRouterPeer, 0, len(e.peers)); for _, p := range e.peers { peers = append(peers, p) }
    routes := make([]FTNCoreRouteState, 0, len(e.routes)); for _, r := range e.routes { routes = append(routes, r) }
    sort.Slice(peers, func(i,j int) bool { return peers[i].ID < peers[j].ID })
    sort.Slice(routes, func(i,j int) bool { return routeKey(routes[i].VRF,routes[i].Prefix) < routeKey(routes[j].VRF,routes[j].Prefix) })
    return FTNCoreRouterSnapshot{Node:e.node, Peers:peers, Routes:routes, CapturedAt:time.Now().UTC()}
}

func (e *FTNCoreRouterEngine) Ready() bool {
    e.mu.RLock(); nodeOK := e.node.Enabled; e.mu.RUnlock()
    return nodeOK && e.bgp.Established()
}
