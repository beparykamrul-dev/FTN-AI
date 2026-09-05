package main

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	api "github.com/osrg/gobgp/v4/api"
)

// FTNGoBGPPeerWatcher polls the live GoBGP peer table and forwards only
// normalized state transitions to the FTN Routed event bridge.
// It never changes peer configuration or routes.
type FTNGoBGPPeerWatcher struct {
	server *FTNGoBGPServer
	bridge *FTNRoutedEventBridge
	interval time.Duration
}

func NewFTNGoBGPPeerWatcher(server *FTNGoBGPServer, bridge *FTNRoutedEventBridge, interval time.Duration) *FTNGoBGPPeerWatcher {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	return &FTNGoBGPPeerWatcher{server: server, bridge: bridge, interval: interval}
}

func (w *FTNGoBGPPeerWatcher) Poll(ctx context.Context) error {
	if w == nil || w.server == nil || w.bridge == nil {
		return fmt.Errorf("GoBGP server and event bridge are required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	w.server.mu.RLock()
	started := w.server.started
	gobgp := w.server.server
	w.server.mu.RUnlock()
	if !started || gobgp == nil {
		return fmt.Errorf("GoBGP server is not started")
	}

	var firstErr error
	err := gobgp.ListPeer(ctx, &api.ListPeerRequest{}, func(peer *api.Peer) {
		if peer == nil || peer.Conf == nil {
			return
		}
		address := peer.Conf.NeighborAddress
		asn := peer.Conf.PeerAsn
		if address == "" || asn == 0 {
			return
		}
		if _, err := netip.ParseAddr(address); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("invalid GoBGP peer address: %w", err)
			}
			return
		}
		established := peer.State != nil && peer.State.SessionState == api.PeerState_ESTABLISHED
		if err := w.bridge.PeerState(ctx, address, asn, established); err != nil && firstErr == nil {
			firstErr = err
		}
	})
	if err != nil {
		return err
	}
	return firstErr
}

func (w *FTNGoBGPPeerWatcher) Run(ctx context.Context) error {
	if w == nil {
		return fmt.Errorf("watcher is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	if err := w.Poll(ctx); err != nil && ctx.Err() == nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.Poll(ctx); err != nil && ctx.Err() == nil {
				return err
			}
		}
	}
}
