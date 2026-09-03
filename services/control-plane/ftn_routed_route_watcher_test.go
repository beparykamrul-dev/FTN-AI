package main

import (
	"testing"
	"time"
)

func TestFTNGoBGPRouteWatcherDefaults(t *testing.T) {
	w := NewFTNGoBGPRouteWatcher(nil, nil, 0)
	if w.interval != 5*time.Second { t.Fatalf("interval=%s", w.interval) }
	if w.last == nil { t.Fatal("route state map must be initialized") }
}

func TestFTNGoBGPRouteWatcherRejectsMissingDependencies(t *testing.T) {
	w := NewFTNGoBGPRouteWatcher(nil, nil, time.Second)
	if err := w.Poll(t.Context()); err == nil { t.Fatal("expected dependency error") }
}
