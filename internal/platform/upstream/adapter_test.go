package upstream

import (
	"testing"
	"time"
)

func TestValidateProfile(t *testing.T) {
	base := Profile{
		Name:            "primary",
		Enabled:         true,
		ASN:             65001,
		RemoteAddresses: []string{"192.0.2.1"},
		RemoteASNs:      []uint32{65002},
		MaxPrefixCount:  100,
	}
	if err := ValidateProfile(base); err != nil {
		t.Fatalf("valid profile rejected: %v", err)
	}

	base.RemoteASNs = nil
	if err := ValidateProfile(base); err == nil {
		t.Fatal("mismatched BGP peer arrays were accepted")
	}
}

func TestRankHealthyEndpoints(t *testing.T) {
	endpoints := []Endpoint{
		{Name: "slow", Address: "192.0.2.2"},
		{Name: "fast", Address: "192.0.2.1"},
		{Name: "down", Address: "192.0.2.3"},
	}
	health := map[string]Health{
		"slow": {Healthy: true, Latency: 40 * time.Millisecond, PacketLossPct: 0.1},
		"fast": {Healthy: true, Latency: 10 * time.Millisecond, PacketLossPct: 0.2},
		"down": {Healthy: false},
	}
	got := RankHealthyEndpoints(endpoints, health)
	if len(got) != 2 || got[0].Name != "fast" || got[1].Name != "slow" {
		t.Fatalf("unexpected ranking: %#v", got)
	}
}
