package main

import (
	"errors"
	"sync"
)

type FTNBGPPeerState struct {
	PeerID string `json:"peer_id"`
	RemoteAddress string `json:"remote_address"`
	RemoteAS uint32 `json:"remote_as"`
	Established bool `json:"established"`
	PrefixesReceived uint64 `json:"prefixes_received"`
	PrefixesAdvertised uint64 `json:"prefixes_advertised"`
	LastError string `json:"last_error,omitempty"`
}

type FTNBGPRouteEvent struct {
	PeerID string `json:"peer_id"`
	Route FTNRoute `json:"route"`
	Direction string `json:"direction"`
}

type FTNGoBGPAdapter struct {
	mu sync.RWMutex
	peers map[string]FTNBGPPeerState
	routes map[string]FTNRoute
	connected bool
}

func NewFTNGoBGPAdapter() *FTNGoBGPAdapter {
	return &FTNGoBGPAdapter{peers: make(map[string]FTNBGPPeerState), routes: make(map[string]FTNRoute)}
}

func (a *FTNGoBGPAdapter) Name() string { return "gobgp-v4" }
func (a *FTNGoBGPAdapter) Established() bool {
	a.mu.RLock(); defer a.mu.RUnlock(); return a.connected
}

// SetSessionState is fed by the real GoBGP server adapter when wired into a live daemon.
func (a *FTNGoBGPAdapter) SetSessionState(connected bool) {
	a.mu.Lock(); defer a.mu.Unlock(); a.connected = connected
}

func (a *FTNGoBGPAdapter) UpsertPeer(state FTNBGPPeerState) error {
	if state.PeerID == "" || state.RemoteAddress == "" || state.RemoteAS == 0 { return errors.New("invalid BGP peer") }
	a.mu.Lock(); defer a.mu.Unlock(); a.peers[state.PeerID] = state; return nil
}

func (a *FTNGoBGPAdapter) Peers() []FTNBGPPeerState {
	a.mu.RLock(); defer a.mu.RUnlock()
	out := make([]FTNBGPPeerState, 0, len(a.peers)); for _, p := range a.peers { out = append(out, p) }; return out
}

func (a *FTNGoBGPAdapter) ApplyRoute(r FTNRoute) error {
	if !a.Established() { return errors.New("BGP session is not established") }
	key := r.VRF + "|" + r.Prefix
	a.mu.Lock(); defer a.mu.Unlock(); a.routes[key] = r; return nil
}

func (a *FTNGoBGPAdapter) WithdrawRoute(r FTNRoute) error {
	key := r.VRF + "|" + r.Prefix
	a.mu.Lock(); defer a.mu.Unlock(); delete(a.routes, key); return nil
}

func (a *FTNGoBGPAdapter) Routes() []FTNRoute {
	a.mu.RLock(); defer a.mu.RUnlock()
	out := make([]FTNRoute, 0, len(a.routes)); for _, r := range a.routes { out = append(out, r) }; return out
}
