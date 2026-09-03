package main

import (
	"context"
	"testing"
)

type bgpTestSink struct{ events []FTNBGPEvent }

func (s *bgpTestSink) Publish(_ context.Context, event FTNBGPEvent) error {
	s.events = append(s.events, event)
	return nil
}

func TestFTNBGPWatcherPublishesStateTransitions(t *testing.T) {
	sink := &bgpTestSink{}
	watcher := NewFTNBGPWatcher(sink)
	ctx := context.Background()
	if err := watcher.UpdatePeer(ctx, "192.0.2.2", 64512, true); err != nil {
		t.Fatal(err)
	}
	if err := watcher.UpdatePeer(ctx, "192.0.2.2", 64512, true); err != nil {
		t.Fatal(err)
	}
	if err := watcher.UpdatePeer(ctx, "192.0.2.2", 64512, false); err != nil {
		t.Fatal(err)
	}
	if len(sink.events) != 2 {
		t.Fatalf("events=%d want 2", len(sink.events))
	}
	if sink.events[0].State != "established" || sink.events[1].State != "idle" {
		t.Fatalf("unexpected events: %+v", sink.events)
	}
}

func TestApplyBFDStateControlsRoutingEligibility(t *testing.T) {
	if !ApplyBFDState(FTNBFDSession{Peer: "192.0.2.2", State: FTNBFDStateUp}, true) {
		t.Fatal("up BFD + established BGP should be eligible")
	}
	if ApplyBFDState(FTNBFDSession{Peer: "192.0.2.2", State: FTNBFDStateDown}, true) {
		t.Fatal("down BFD must not be eligible")
	}
	if ApplyBFDState(FTNBFDSession{Peer: "192.0.2.2", State: FTNBFDStateUp}, false) {
		t.Fatal("down BGP must not be eligible")
	}
}

func TestFTNBGPWatcherRejectsInvalidPeer(t *testing.T) {
	watcher := NewFTNBGPWatcher(nil)
	if err := watcher.UpdatePeer(context.Background(), "", 64512, true); err == nil {
		t.Fatal("expected invalid peer error")
	}
	if err := watcher.UpdatePeer(context.Background(), "192.0.2.2", 0, true); err == nil {
		t.Fatal("expected invalid ASN error")
	}
}
