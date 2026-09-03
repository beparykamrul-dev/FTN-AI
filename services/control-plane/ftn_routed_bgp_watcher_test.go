package main

import (
	"context"
	"testing"
	"time"
)

func TestFTNGoBGPPeerWatcherRequiresDependencies(t *testing.T) {
	w := NewFTNGoBGPPeerWatcher(nil, nil, time.Second)
	if err := w.Poll(context.Background()); err == nil {
		t.Fatal("expected dependency error")
	}
}

func TestFTNGoBGPPeerWatcherDefaultInterval(t *testing.T) {
	w := NewFTNGoBGPPeerWatcher(nil, nil, 0)
	if w.interval != 5*time.Second {
		t.Fatalf("interval=%s want 5s", w.interval)
	}
}
