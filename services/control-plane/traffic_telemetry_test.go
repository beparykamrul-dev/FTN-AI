package main

import (
	"testing"
	"time"
)

func TestTrafficQualityStoreRejectsInvalidMetrics(t *testing.T) {
	s := NewTrafficQualityStore()
	now := time.Now().UTC()
	if err := s.Upsert(TrafficQualityObservation{ServiceID:"pubg", PathID:"edge-a", PacketLoss:101, ObservedAt:now}, now); err == nil {
		t.Fatal("expected packet-loss validation failure")
	}
	if err := s.Upsert(TrafficQualityObservation{ServiceID:"pubg", PathID:"edge-a", Congestion:1.1, ObservedAt:now}, now); err == nil {
		t.Fatal("expected congestion validation failure")
	}
}

func TestTrafficQualityStoreRejectsStaleObservation(t *testing.T) {
	s := NewTrafficQualityStore()
	now := time.Now().UTC()
	if err := s.Upsert(TrafficQualityObservation{ServiceID:"whatsapp", PathID:"edge-a", ObservedAt:now.Add(-3*time.Minute)}, now); err == nil {
		t.Fatal("expected stale observation rejection")
	}
}

func TestTrafficQualityStoreSnapshotExpiresOldEntries(t *testing.T) {
	s := NewTrafficQualityStore()
	now := time.Now().UTC()
	if err := s.Upsert(TrafficQualityObservation{ServiceID:"telegram", PathID:"edge-a", LatencyMs:20, ObservedAt:now}, now); err != nil {
		t.Fatal(err)
	}
	if got := s.Snapshot(now.Add(3*time.Minute)); len(got) != 0 {
		t.Fatalf("expected expired store to be empty, got %d", len(got))
	}
}
