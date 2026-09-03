package main

import (
	"context"
	"testing"
)

type routedEventSink struct{ events []FTNBGPEvent }

func (s *routedEventSink) Publish(_ context.Context, event FTNBGPEvent) error {
	s.events = append(s.events, event)
	return nil
}

func TestFTNRoutedEventBridgeEligibility(t *testing.T) {
	sink := &routedEventSink{}
	b := NewFTNRoutedEventBridge(sink)
	ctx := context.Background()
	if err := b.PeerState(ctx, "192.0.2.2", 64512, true); err != nil { t.Fatal(err) }
	if err := b.BFDState(ctx, "192.0.2.2", FTNBFDDown); err != nil { t.Fatal(err) }
	if b.Eligible("192.0.2.2") { t.Fatal("down BFD session must not be eligible") }
	if err := b.BFDState(ctx, "192.0.2.2", FTNBFDUp); err != nil { t.Fatal(err) }
	if !b.Eligible("192.0.2.2") { t.Fatal("established BGP + BFD up must be eligible") }
	if len(sink.events) != 3 { t.Fatalf("events=%d want 3", len(sink.events)) }
}

func TestFTNRoutedEventBridgeSuppressesDuplicateState(t *testing.T) {
	sink := &routedEventSink{}
	b := NewFTNRoutedEventBridge(sink)
	ctx := context.Background()
	if err := b.PeerState(ctx, "198.51.100.2", 64513, true); err != nil { t.Fatal(err) }
	if err := b.PeerState(ctx, "198.51.100.2", 64513, true); err != nil { t.Fatal(err) }
	if len(sink.events) != 1 { t.Fatalf("events=%d want 1", len(sink.events)) }
}

func TestFTNRoutedEventBridgeNormalizesRouteEvents(t *testing.T) {
	sink := &routedEventSink{}
	b := NewFTNRoutedEventBridge(sink)
	if err := b.RouteEvent(context.Background(), "203.0.113.7/24", true); err != nil { t.Fatal(err) }
	if len(sink.events) != 1 || sink.events[0].Prefix != "203.0.113.0/24" || sink.events[0].Type != FTNBGPEventRouteAdd {
		t.Fatalf("unexpected event: %+v", sink.events)
	}
}
