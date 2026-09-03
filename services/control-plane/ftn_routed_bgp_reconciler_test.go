package main

import (
	"context"
	"testing"
)

type fakeBGPPeerSource struct { peers []FTNObservedBGPPeer }
func (f *fakeBGPPeerSource) ListPeers(context.Context) ([]FTNObservedBGPPeer, error) { return append([]FTNObservedBGPPeer(nil), f.peers...), nil }

func TestFTNBGPPeerReconcilerReportsPeerDisappearance(t *testing.T) {
	source := &fakeBGPPeerSource{peers: []FTNObservedBGPPeer{{Address:"192.0.2.2", ASN:64512, Established:true}}}
	sink := &routedEventSink{}
	r := NewFTNBGPPeerReconciler(source, NewFTNRoutedEventBridge(sink))
	ctx := context.Background()
	if err := r.Reconcile(ctx); err != nil { t.Fatal(err) }
	if len(sink.events) != 1 || sink.events[0].State != "established" { t.Fatalf("events=%+v", sink.events) }
	source.peers = nil
	if err := r.Reconcile(ctx); err != nil { t.Fatal(err) }
	if len(sink.events) != 2 || sink.events[1].State != "idle" { t.Fatalf("events=%+v", sink.events) }
}

func TestFTNBGPPeerReconcilerIgnoresInvalidPeers(t *testing.T) {
	source := &fakeBGPPeerSource{peers: []FTNObservedBGPPeer{{Address:"bad", ASN:64512, Established:true},{Address:"192.0.2.2", ASN:0, Established:true}}}
	sink := &routedEventSink{}
	r := NewFTNBGPPeerReconciler(source, NewFTNRoutedEventBridge(sink))
	if err := r.Reconcile(context.Background()); err != nil { t.Fatal(err) }
	if len(sink.events) != 0 { t.Fatalf("events=%+v", sink.events) }
}
