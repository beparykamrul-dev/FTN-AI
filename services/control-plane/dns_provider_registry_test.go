package main

import "testing"

func TestRankDNSProviders(t *testing.T) {
	health := []DNSProviderHealth{
		{ProviderID: "slow", Endpoint: "slow:53", Healthy: true, LatencyMS: 80},
		{ProviderID: "fast", Endpoint: "fast:53", Healthy: true, LatencyMS: 12},
		{ProviderID: "failed", Endpoint: "failed:53", Healthy: false, LatencyMS: 1},
		{ProviderID: "too-slow", Endpoint: "too-slow:53", Healthy: true, LatencyMS: 700},
	}
	got := RankDNSProviders(health, 500)
	if len(got) != 2 || got[0].ProviderID != "fast" || got[1].ProviderID != "slow" {
		t.Fatalf("unexpected ranking: %+v", got)
	}
}

func TestPreferredDNSEndpoints(t *testing.T) {
	p := DNSProvider{ID: "x", Enabled: true, Endpoints: []DNSProviderEndpoint{
		{Protocol: "https", Address: "https://x/dns-query"},
		{Protocol: "udp", Address: "x:53"},
		{Protocol: "tls", Address: "x:853"},
	}}
	got := PreferredDNSEndpoints(p, []string{"udp", "tls"})
	if len(got) != 2 || got[0].Protocol != "udp" || got[1].Protocol != "tls" {
		t.Fatalf("unexpected endpoints: %+v", got)
	}
}
