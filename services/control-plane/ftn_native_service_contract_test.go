package main

import "testing"

func TestFTNNativeServiceProtocolContract(t *testing.T) {
	endpoints := []FTNNativeProtocolEndpoint{
		{Service: "ftn-dns", Protocol: "doq", Transport: "udp", Healthy: true, Authorized: true, LatencyMS: 8, LossPct: 0.1, Capacity: 95},
		{Service: "ftn-dns", Protocol: "doh", Transport: "tcp", Healthy: true, Authorized: true, LatencyMS: 12, LossPct: 0.1, Capacity: 90},
		{Service: "ftn-dns", Protocol: "dns", Transport: "udp", Healthy: false, Authorized: true, LatencyMS: 1, Capacity: 100},
	}
	got, ok := NegotiateFTNNativeProtocol("ftn-dns", []string{"dns", "doh", "doq"}, []string{"udp", "tcp"}, endpoints)
	if !ok || got.Protocol != "doq" { t.Fatalf("expected healthy lowest-latency FTN-native endpoint, got %+v", got) }
}
