package main

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"
)

// FTNGoBGPRouteEventSource converts normalized GoBGP RIB snapshots into
// bounded routing events. It never applies, withdraws, or mutates routes.
type FTNGoBGPRouteEventSource struct {
	server *FTNGoBGPServer
	bridge *FTNRoutedEventBridge
}

func NewFTNGoBGPRouteEventSource(server *FTNGoBGPServer, bridge *FTNRoutedEventBridge) *FTNGoBGPRouteEventSource {
	return &FTNGoBGPRouteEventSource{server: server, bridge: bridge}
}

func (s *FTNGoBGPRouteEventSource) PublishSnapshot(ctx context.Context, routes []FTNRoute) error {
	if s == nil || s.bridge == nil {
		return fmt.Errorf("event bridge is required")
	}
	for _, route := range routes {
		prefix := strings.TrimSpace(route.Prefix)
		if _, err := netip.ParsePrefix(prefix); err != nil {
			continue
		}
		if err := s.bridge.RouteEvent(ctx, prefix, route.Active); err != nil {
			return err
		}
	}
	return nil
}

// PollOnce reads the live local RIB and emits normalized route-state events.
// A caller owns the polling schedule so this boundary can be integrated into
// an existing worker lifecycle without creating an uncontrolled goroutine.
func (s *FTNGoBGPRouteEventSource) PollOnce(ctx context.Context) error {
	if s == nil || s.server == nil {
		return fmt.Errorf("GoBGP server is required")
	}
	if s.bridge == nil {
		return fmt.Errorf("event bridge is required")
	}
	routes, err := s.server.ListIPv4Routes(ctx)
	if err != nil {
		return err
	}
	return s.PublishSnapshot(ctx, routes)
}

// Run periodically samples the local GoBGP RIB until context cancellation.
// It is telemetry-only; no route mutation is performed here.
func (s *FTNGoBGPRouteEventSource) Run(ctx context.Context, interval time.Duration) error {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	if err := s.PollOnce(ctx); err != nil && ctx.Err() == nil {
		return err
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.PollOnce(ctx); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return err
			}
		}
	}
}
