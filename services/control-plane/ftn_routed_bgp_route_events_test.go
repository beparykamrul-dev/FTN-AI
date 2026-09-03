package main

import (
	"context"
	"testing"
)

func TestFTNGoBGPRouteEventSourcePublishSnapshot(t *testing.T) {
	sink := &routedEventSink{}
	bridge := NewFTNRoutedEventBridge(sink)
	source := NewFTNGoBGPRouteEventSource(nil, bridge)

	routes := []FTNRoute{
		{Prefix: "203.0.113.7/24", Protocol: "bgp", Active: true},
		{Prefix: "198.51.100.9/24", Protocol: "bgp", Active: false},
		{Prefix: "not-a-prefix", Protocol: "bgp", Active: true},
	}
	if err := source.PublishSnapshot(context.Background(), routes); err != nil {
		t.Fatal(err)
	}
	if len(sink.events) != 2 {
		t.Fatalf("events=%d want 2", len(sink.events))
	}
	if sink.events[0].Prefix != "203.0.113.0/24" || sink.events[0].Type != FTNBGPEventRouteAdd {
		t.Fatalf("unexpected add event: %+v", sink.events[0])
	}
	if sink.events[1].Prefix != "198.51.100.0/24" || sink.events[1].Type != FTNBGPEventRouteDel {
		t.Fatalf("unexpected delete event: %+v", sink.events[1])
	}
}

func TestFTNGoBGPRouteEventSourceRequiresBridge(t *testing.T) {
	source := NewFTNGoBGPRouteEventSource(nil, nil)
	if err := source.PublishSnapshot(context.Background(), nil); err == nil {
		t.Fatal("expected missing bridge error")
	}
}
