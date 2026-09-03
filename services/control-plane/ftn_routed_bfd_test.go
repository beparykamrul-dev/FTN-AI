package main

import (
	"context"
	"testing"
)

type fakeBFDSource struct { observations []FTNBFDObservation }
func (f *fakeBFDSource) ListBFD(context.Context) ([]FTNBFDObservation, error) { return append([]FTNBFDObservation(nil), f.observations...), nil }

func TestFTNBFDReconcilerPublishesState(t *testing.T) {
	source := &fakeBFDSource{observations: []FTNBFDObservation{{Peer: "192.0.2.2", State: FTNBFDUp}}}
	sink := &routedEventSink{}
	r := NewFTNBFDReconciler(source, NewFTNRoutedEventBridge(sink))
	if err := r.Reconcile(context.Background()); err != nil { t.Fatal(err) }
	if len(sink.events) != 1 || sink.events[0].State != "bfd-up" { t.Fatalf("events=%+v", sink.events) }
}

func TestFTNBFDReconcilerMarksMissingPeerUnknown(t *testing.T) {
	source := &fakeBFDSource{observations: []FTNBFDObservation{{Peer: "192.0.2.2", State: FTNBFDUp}}}
	sink := &routedEventSink{}
	r := NewFTNBFDReconciler(source, NewFTNRoutedEventBridge(sink))
	ctx := context.Background()
	if err := r.Reconcile(ctx); err != nil { t.Fatal(err) }
	source.observations = nil
	if err := r.Reconcile(ctx); err != nil { t.Fatal(err) }
	if len(sink.events) != 2 || sink.events[1].State != "bfd-unknown" { t.Fatalf("events=%+v", sink.events) }
}

func TestFTNBFDReconcilerIgnoresInvalidObservation(t *testing.T) {
	source := &fakeBFDSource{observations: []FTNBFDObservation{{Peer: "bad", State: FTNBFDUp}}}
	sink := &routedEventSink{}
	r := NewFTNBFDReconciler(source, NewFTNRoutedEventBridge(sink))
	if err := r.Reconcile(context.Background()); err != nil { t.Fatal(err) }
	if len(sink.events) != 0 { t.Fatalf("events=%+v", sink.events) }
}
