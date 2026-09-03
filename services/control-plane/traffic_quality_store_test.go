package main

import (
    "testing"
    "time"
)

func TestTrafficQualityStoreRejectsInvalidAndStale(t *testing.T) {
    now := time.Unix(1000, 0)
    store := NewTrafficQualityStore()
    if err := store.Upsert(TrafficQualityObservation{PathID: "p1", ServiceID: "pubg", Healthy: true, LatencyMs: -1, ObservedAt: now}, now); err == nil {
        t.Fatal("expected invalid latency rejection")
    }
    if err := store.Upsert(TrafficQualityObservation{PathID: "p1", ServiceID: "pubg", Healthy: true, ObservedAt: now.Add(-trafficQualityTTL - time.Second)}, now); err == nil {
        t.Fatal("expected stale observation rejection")
    }
}

func TestTrafficQualityStoreSnapshotExpires(t *testing.T) {
    now := time.Unix(1000, 0)
    store := NewTrafficQualityStore()
    if err := store.Upsert(TrafficQualityObservation{PathID: "p1", ServiceID: "pubg", Class: TrafficGaming, LatencyMs: 12, JitterMs: 2, PacketLoss: 0.1, Congestion: 0.2, Healthy: true, ObservedAt: now}, now); err != nil {
        t.Fatal(err)
    }
    if got := store.Snapshot("pubg", now.Add(30*time.Second)); len(got) != 1 || got[0].PathID != "p1" {
        t.Fatalf("unexpected snapshot: %+v", got)
    }
    if got := store.Snapshot("pubg", now.Add(trafficQualityTTL+time.Second)); len(got) != 0 {
        t.Fatalf("expected expired snapshot, got %+v", got)
    }
    if removed := store.Prune(now.Add(trafficQualityTTL + time.Second)); removed != 1 {
        t.Fatalf("expected one prune, got %d", removed)
    }
}

func TestTrafficRuntimeFlowCountIsRaceSafe(t *testing.T) {
    runtime := NewTrafficRuntime()
    if err := runtime.UpsertEndpoint(ManagedEndpoint{ServiceID: "pubg", CIDR: "198.51.100.0/24"}); err != nil {
        t.Fatal(err)
    }
    now := time.Unix(1000, 0)
    if accepted := runtime.Ingest([]FlowRecord{{SourceIP: "192.0.2.1", DestinationIP: "198.51.100.1"}}, now); accepted != 1 {
        t.Fatalf("accepted=%d", accepted)
    }
    if got := runtime.FlowCount(); got != 1 {
        t.Fatalf("flow count=%d", got)
    }
}
