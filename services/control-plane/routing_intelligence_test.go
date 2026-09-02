package main

import "testing"

func TestNormalizeRouteIntelligenceInput(t *testing.T) {
	in, ok := NormalizeRouteIntelligenceInput(RouteIntelligenceInput{Prefix: "2001:db8::/32", Country: "bd", LatencyClass: "LOW", TrafficClass: "NORMAL", PeerRole: "EDGE", RouteSource: "EXTERNAL"})
	if !ok || in.Prefix != "2001:db8::/32" || in.Country != "BD" || in.LatencyClass != "low" || in.PeerRole != "edge" {
		t.Fatalf("unexpected normalization: %#v ok=%v", in, ok)
	}
}

func TestBuildRouteDecisionAlwaysRequiresApproval(t *testing.T) {
	d := BuildRouteDecision(RouteIntelligenceInput{Prefix: "203.0.113.0/24", ASN: 64512, RouteSource: "internal"})
	if !d.RequiresApproval {
		t.Fatal("route mutation must remain approval-gated")
	}
	if d.Decision != "allow" || d.LocalPref != 100 || d.DecisionHash == "" {
		t.Fatalf("unexpected decision: %#v", d)
	}
}

func TestBuildRouteDecisionRejectsUntrustedExternalRoute(t *testing.T) {
	d := BuildRouteDecision(RouteIntelligenceInput{Prefix: "198.51.100.0/24", RouteSource: "external", TrafficClass: "untrusted"})
	if d.Decision != "deny" || d.Reason != "untrusted_external_route" || !d.RequiresApproval {
		t.Fatalf("unexpected decision: %#v", d)
	}
}
