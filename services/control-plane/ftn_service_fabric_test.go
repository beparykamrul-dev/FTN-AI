package main

import "testing"

func TestNormalizeFTNServiceEndpoint(t *testing.T) {
	e, err := NormalizeFTNServiceEndpoint(FTNServiceEndpoint{
		Service: " FTN-API ", Protocol: "HTTP3", Transport: "UDP", Address: "10.0.0.10:443",
		Healthy: true, Capacity: 100, FTNOwned: true,
	})
	if err != nil { t.Fatal(err) }
	if e.Service != "ftn-api" || e.Protocol != "http3" || e.Transport != "udp" { t.Fatalf("unexpected normalization: %+v", e) }
}

func TestNormalizeFTNServiceEndpointRejectsUnauthorized(t *testing.T) {
	_, err := NormalizeFTNServiceEndpoint(FTNServiceEndpoint{Service: "dns", Protocol: "doq", Transport: "udp", Address: "1.2.3.4:853"})
	if err == nil { t.Fatal("expected authorization error") }
}

func TestSelectFTNServiceEndpointsDeterministic(t *testing.T) {
	got := SelectFTNServiceEndpoints("ftn-api", []FTNServiceEndpoint{
		{Service: "ftn-api", Protocol: "http3", Transport: "udp", Address: "b", Healthy: true, LatencyMS: 10, Capacity: 80, FTNOwned: true},
		{Service: "ftn-api", Protocol: "http2", Transport: "tcp", Address: "a", Healthy: true, LatencyMS: 10, Capacity: 80, FTNOwned: true},
		{Service: "ftn-api", Protocol: "http1", Transport: "tcp", Address: "c", Healthy: false, LatencyMS: 1, FTNOwned: true},
	}, 2)
	if len(got) != 2 || got[0].Address != "a" || got[1].Address != "b" { t.Fatalf("unexpected selection: %+v", got) }
}
