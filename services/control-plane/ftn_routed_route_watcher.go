package main

import (
	"context"
	"fmt"
	"net/netip"
	"sync"
	"time"

	api "github.com/osrg/gobgp/v4/api"
)

// FTNGoBGPRouteWatcher converts changes in the live GoBGP local RIB into
// normalized FTN Routed events. It is telemetry-only: it never changes routes.
type FTNGoBGPRouteWatcher struct {
	server *FTNGoBGPServer
	bridge *FTNRoutedEventBridge
	interval time.Duration
	mu sync.Mutex
	last map[string]bool
}

func NewFTNGoBGPRouteWatcher(server *FTNGoBGPServer, bridge *FTNRoutedEventBridge, interval time.Duration) *FTNGoBGPRouteWatcher {
	if interval <= 0 { interval = 5 * time.Second }
	return &FTNGoBGPRouteWatcher{server: server, bridge: bridge, interval: interval, last: make(map[string]bool)}
}

func (w *FTNGoBGPRouteWatcher) Poll(ctx context.Context) error {
	if w == nil || w.server == nil || w.bridge == nil { return fmt.Errorf("GoBGP server and event bridge are required") }
	w.server.mu.RLock(); started, gobgp := w.server.started, w.server.server; w.server.mu.RUnlock()
	if !started || gobgp == nil { return fmt.Errorf("GoBGP server is not started") }
	current := make(map[string]bool)
	family := &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST}
	err := gobgp.ListPath(ctx, &api.ListPathRequest{TableType: api.TableType_GLOBAL, Family: family}, func(destination *api.Destination) {
		if destination == nil { return }
		for _, path := range destination.Paths {
			if path == nil || path.Nlri == nil || path.IsWithdraw { continue }
			pfx := path.GetNlri().GetPrefix(); if pfx == nil { continue }
			p, err := netip.ParsePrefix(fmt.Sprintf("%s/%d", pfx.Prefix, pfx.PrefixLen)); if err != nil { continue }
			current[p.Masked().String()] = true
		}
	})
	if err != nil { return err }
	w.mu.Lock(); previous := w.last; w.last = current; w.mu.Unlock()
	for prefix := range current { if !previous[prefix] { if err := w.bridge.RouteEvent(ctx, prefix, true); err != nil { return err } } }
	for prefix := range previous { if !current[prefix] { if err := w.bridge.RouteEvent(ctx, prefix, false); err != nil { return err } } }
	return nil
}

func (w *FTNGoBGPRouteWatcher) Run(ctx context.Context) error {
	if w == nil { return fmt.Errorf("route watcher is required") }
	if err := w.Poll(ctx); err != nil && ctx.Err() == nil { return err }
	ticker := time.NewTicker(w.interval); defer ticker.Stop()
	for { select {
	case <-ctx.Done(): return ctx.Err()
	case <-ticker.C: if err := w.Poll(ctx); err != nil && ctx.Err() == nil { return err }
	} }
}
