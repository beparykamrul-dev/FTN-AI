package main

import (
	"context"
	"fmt"
	"net/netip"
	"sort"
	"sync"
	"time"
)

type FTNBGPEventType string

const (
	FTNBGPEventPeerState FTNBGPEventType = "bgp.peer.state"
	FTNBGPEventRouteAdd  FTNBGPEventType = "bgp.route.add"
	FTNBGPEventRouteDel  FTNBGPEventType = "bgp.route.delete"
)

type FTNBGPEvent struct {
	Type      FTNBGPEventType `json:"type"`
	Peer      string          `json:"peer,omitempty"`
	Prefix    string          `json:"prefix,omitempty"`
	State     string          `json:"state,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type FTNBGPEventSink interface {
	Publish(context.Context, FTNBGPEvent) error
}

type FTNBGPState struct {
	Peer        string    `json:"peer"`
	ASN         uint32    `json:"asn"`
	Established bool      `json:"established"`
	LastChange  time.Time `json:"last_change"`
}

type FTNBGPWatcher struct {
	mu    sync.RWMutex
	peers map[string]FTNBGPState
	sink  FTNBGPEventSink
}

func NewFTNBGPWatcher(sink FTNBGPEventSink) *FTNBGPWatcher {
	return &FTNBGPWatcher{peers: make(map[string]FTNBGPState), sink: sink}
}

func (w *FTNBGPWatcher) UpdatePeer(ctx context.Context, peer string, asn uint32, established bool) error {
	if w == nil || ctx == nil {
		return fmt.Errorf("watcher and context are required")
	}
	if peer == "" || asn == 0 {
		return fmt.Errorf("peer and ASN are required")
	}
	if _, err := netip.ParseAddr(peer); err != nil {
		return fmt.Errorf("invalid peer address: %w", err)
	}
	w.mu.Lock()
	previous, exists := w.peers[peer]
	state := FTNBGPState{Peer: peer, ASN: asn, Established: established, LastChange: time.Now().UTC()}
	if exists && previous.Established == established && previous.ASN == asn {
		w.mu.Unlock()
		return nil
	}
	w.peers[peer] = state
	w.mu.Unlock()
	if w.sink != nil {
		status := "idle"
		if established {
			status = "established"
		}
		return w.sink.Publish(ctx, FTNBGPEvent{Type: FTNBGPEventPeerState, Peer: peer, State: status, Timestamp: state.LastChange})
	}
	return nil
}

func (w *FTNBGPWatcher) State(peer string) (FTNBGPState, bool) {
	if w == nil {
		return FTNBGPState{}, false
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	s, ok := w.peers[peer]
	return s, ok
}

func (w *FTNBGPWatcher) Snapshot() []FTNBGPState {
	if w == nil {
		return nil
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]FTNBGPState, 0, len(w.peers))
	for _, state := range w.peers {
		out = append(out, state)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Peer != out[j].Peer {
			return out[i].Peer < out[j].Peer
		}
		return out[i].ASN < out[j].ASN
	})
	return out
}

type FTNBFDState string

const (
	FTNBFDStateUp      FTNBFDState = "up"
	FTNBFDStateDown    FTNBFDState = "down"
	FTNBFDStateUnknown FTNBFDState = "unknown"
)

type FTNBFDSession struct {
	Peer       string      `json:"peer"`
	State      FTNBFDState `json:"state"`
	LastChange time.Time   `json:"last_change"`
}

func ApplyBFDState(session FTNBFDSession, bgpEstablished bool) bool {
	return bgpEstablished && session.State == FTNBFDStateUp
}
