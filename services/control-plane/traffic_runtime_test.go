package main

import (
    "testing"
    "time"
)

func TestTrafficRuntimeClassifiesManagedEndpoint(t *testing.T) {
    rt := NewTrafficRuntime()
    if err := rt.UpsertEndpoint(ManagedEndpoint{ServiceID:"pubg", CIDR:"203.0.113.0/24"}); err != nil { t.Fatal(err) }
    obs, ok := rt.Classify(FlowRecord{SourceIP:"10.0.0.2", DestinationIP:"203.0.113.10", DestinationPort:443, Bytes:1000}, time.Unix(1000,0))
    if !ok { t.Fatal("expected managed endpoint match") }
    if obs.ServiceID != "pubg" || obs.Class != TrafficGaming { t.Fatalf("unexpected classification: %+v", obs) }
}

func TestTrafficRuntimeUsesLongestPrefix(t *testing.T) {
    rt := NewTrafficRuntime()
    _ = rt.UpsertEndpoint(ManagedEndpoint{ServiceID:"whatsapp", CIDR:"198.51.100.0/24"})
    _ = rt.UpsertEndpoint(ManagedEndpoint{ServiceID:"telegram", CIDR:"198.51.100.0/25"})
    obs, ok := rt.Classify(FlowRecord{DestinationIP:"198.51.100.20"}, time.Unix(1000,0))
    if !ok || obs.ServiceID != "telegram" { t.Fatalf("expected longest-prefix service, got %+v", obs) }
}

func TestTrafficRuntimeExpiresEndpoint(t *testing.T) {
    rt := NewTrafficRuntime()
    now := time.Unix(1000,0)
    _ = rt.UpsertEndpoint(ManagedEndpoint{ServiceID:"imo", CIDR:"192.0.2.0/24", ExpiresAt:now.Add(-time.Second)})
    if _, ok := rt.Classify(FlowRecord{DestinationIP:"192.0.2.10"}, now); ok { t.Fatal("expired endpoint must not classify") }
}
