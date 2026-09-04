package main

import (
	"context"
	"fmt"
	"net/netip"
	"sync"
	"time"
)

// FTNRoutedEventBridge converts routing-state changes into bounded FTN events.
// It does not execute route changes and does not persist raw protocol payloads.
type FTNRoutedEventBridge struct {
	mu    sync.RWMutex
	sink  FTNBGPEventSink
	bfd   map[string]FTNBFDSession
	peers map[string]FTNBGPState
}

func NewFTNRoutedEventBridge(sink FTNBGPEventSink) *FTNRoutedEventBridge {
	return &FTNRoutedEventBridge{sink: sink, bfd: make(map[string]FTNBFDSession), peers: make(map[string]FTNBGPState)}
}

func (b *FTNRoutedEventBridge) PeerState(ctx context.Context, peer string, asn uint32, established bool) error {
	if b == nil || ctx == nil {
		return fmt.Errorf("bridge and context are required")
	}
	if peer == "" || asn == 0 {
		return fmt.Errorf("peer and ASN are required")
	}
	if _, err := netip.ParseAddr(peer); err != nil {
		return fmt.Errorf("invalid peer address: %w", err)
	}
	now := time.Now().UTC()
	b.mu.Lock()
	previous, exists := b.peers[peer]
	if exists && previous.ASN == asn && previous.Established == established {
		b.mu.Unlock()
		return nil
	}
	b.peers[peer] = FTNBGPState{Peer: peer, ASN: asn, Established: established, LastChange: now}
	b.mu.Unlock()
	return b.publish(ctx, FTNBGPEvent{Type: FTNBGPEventPeerState, Peer: peer, State: bgpStateName(established), Timestamp: now})
}

func (b *FTNRoutedEventBridge) BFDState(ctx context.Context, peer string, state FTNBFDState) error {
	if b == nil || ctx == nil {
		return fmt.Errorf("bridge and context are required")
	}
	if peer == "" {
		return fmt.Errorf("peer is required")
	}
	if _, err := netip.ParseAddr(peer); err != nil {
		return fmt.Errorf("invalid peer address: %w", err)
	}
	switch state {
	case FTNBFDUp, FTNBFDDown, FTNBFDUnknown:
	default:
		return fmt.Errorf("invalid BFD state")
	}
	now := time.Now().UTC()
	b.mu.Lock()
	previous, exists := b.bfd[peer]
	if exists && previous.State == state {
		b.mu.Unlock()
		return nil
	}
	b.bfd[peer] = FTNBFDSession{Peer: peer, State: state, LastChange: now}
	b.mu.Unlock()
	return b.publish(ctx, FTNBGPEvent{Type: FTNBGPEventPeerState, Peer: peer, State: "bfd-" + string(state), Timestamp: now})
}

func (b *FTNRoutedEventBridge) RouteEvent(ctx context.Context, prefix string, added bool) error {
	if b == nil || ctx == nil {
		return fmt.Errorf("bridge and context are required")
	}
	p, err := netip.ParsePrefix(prefix)
	if err != nil {
		return fmt.Errorf("invalid prefix: %w", err)
	}
	p = p.Masked()
	eventType := FTNBGPEventRouteDel
	if added {
		eventType = FTNBGPEventRouteAdd
	}
	return b.publish(ctx, FTNBGPEvent{Type: eventType, Prefix: p.String(), Timestamp: time.Now().UTC()})
}

func (b *FTNRoutedEventBridge) Eligible(peer string) bool {
	if b == nil {
		return false
	}
	if _, err := netip.ParseAddr(peer); err != nil {
		return false
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	bgp, ok := b.peers[peer]
	bfd, bfdOK := b.bfd[peer]
	return ok && bfdOK && bgp.Established && bfd.State == FTNBFDUp
}

func (b *FTNRoutedEventBridge) publish(ctx context.Context, event FTNBGPEvent) error {
	if b == nil || ctx == nil {
		return fmt.Errorf("bridge and context are required")
	}
	if b.sink == nil {
		return nil
	}
	return b.sink.Publish(ctx, event)
}

func bgpStateName(established bool) string {
	if established {
		return "established"
	}
	return "idle"
}
