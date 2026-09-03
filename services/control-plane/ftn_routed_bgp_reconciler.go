package main

import (
	"context"
	"fmt"
	"net/netip"
)

// FTNBGPPeerReconciler detects both state transitions and peer disappearance.
// A missing peer is reported as idle so stale established state cannot survive
// indefinitely in the routing event stream.
type FTNBGPPeerReconciler struct {
	previous map[string]FTNObservedBGPPeer
	bridge   *FTNRoutedEventBridge
	source   FTNBGPPeerSource
}

func NewFTNBGPPeerReconciler(source FTNBGPPeerSource, bridge *FTNRoutedEventBridge) *FTNBGPPeerReconciler {
	return &FTNBGPPeerReconciler{source: source, bridge: bridge, previous: make(map[string]FTNObservedBGPPeer)}
}

func (r *FTNBGPPeerReconciler) Reconcile(ctx context.Context) error {
	if r == nil || r.source == nil || r.bridge == nil {
		return fmt.Errorf("BGP peer source and event bridge are required")
	}
	current, err := r.source.ListPeers(ctx)
	if err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(current))
	for _, peer := range current {
		if _, err := netip.ParseAddr(peer.Address); err != nil || peer.ASN == 0 {
			continue
		}
		seen[peer.Address] = struct{}{}
		if err := r.bridge.PeerState(ctx, peer.Address, peer.ASN, peer.Established); err != nil {
			return err
		}
		r.previous[peer.Address] = peer
	}
	for address, peer := range r.previous {
		if _, ok := seen[address]; ok {
			continue
		}
		if peer.Established {
			if err := r.bridge.PeerState(ctx, address, peer.ASN, false); err != nil {
				return err
			}
		}
		delete(r.previous, address)
	}
	return nil
}
