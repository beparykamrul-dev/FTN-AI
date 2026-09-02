package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"
)

type RouteIntelligenceInput struct {
	Prefix       string `json:"prefix"`
	ASN          uint32 `json:"asn"`
	Country      string `json:"country"`
	LatencyClass string `json:"latency_class"`
	TrafficClass string `json:"traffic_class"`
	PeerRole     string `json:"peer_role"`
	RouteSource  string `json:"route_source"`
}

type RouteIntelligenceDecision struct {
	Decision         string   `json:"decision"`
	LocalPref        uint32   `json:"local_preference"`
	MED              uint32   `json:"med"`
	Communities      []string `json:"communities,omitempty"`
	Reason           string   `json:"reason"`
	DecisionHash     string   `json:"decision_hash"`
	RequiresApproval bool     `json:"requires_approval"`
}

func NormalizeRouteIntelligenceInput(in RouteIntelligenceInput) (RouteIntelligenceInput, bool) {
	p, err := netip.ParsePrefix(strings.TrimSpace(in.Prefix))
	if err != nil || !p.IsValid() {
		return RouteIntelligenceInput{}, false
	}
	in.Prefix = p.String()
	in.Country = strings.ToUpper(strings.TrimSpace(in.Country))
	in.LatencyClass = strings.ToLower(strings.TrimSpace(in.LatencyClass))
	in.TrafficClass = strings.ToLower(strings.TrimSpace(in.TrafficClass))
	in.PeerRole = strings.ToLower(strings.TrimSpace(in.PeerRole))
	in.RouteSource = strings.ToLower(strings.TrimSpace(in.RouteSource))
	return in, true
}

func HashRouteDecision(in RouteIntelligenceInput, d RouteIntelligenceDecision) string {
	b := fmt.Sprintf("%s|%d|%s|%s|%s|%s|%s|%s|%d|%d|%s|%t|%s",
		in.Prefix, in.ASN, in.Country, in.LatencyClass, in.TrafficClass,
		in.PeerRole, in.RouteSource, d.Decision, d.LocalPref, d.MED,
		strings.Join(d.Communities, ","), d.RequiresApproval, d.Reason)
	h := sha256.Sum256([]byte(b))
	return hex.EncodeToString(h[:])
}

// BuildRouteDecision creates an advisory decision. It never changes routing state.
// A separate approval/execution workflow must authorize any route mutation.
func BuildRouteDecision(in RouteIntelligenceInput) RouteIntelligenceDecision {
	in, ok := NormalizeRouteIntelligenceInput(in)
	if !ok {
		d := RouteIntelligenceDecision{Decision: "deny", Reason: "invalid_route_input", RequiresApproval: true}
		d.DecisionHash = HashRouteDecision(RouteIntelligenceInput{}, d)
		return d
	}
	d := RouteIntelligenceDecision{Decision: "allow", LocalPref: 100, Reason: "default_policy", RequiresApproval: true}
	if in.RouteSource == "external" && in.TrafficClass == "untrusted" {
		d.Decision = "deny"
		d.Reason = "untrusted_external_route"
	}
	d.DecisionHash = HashRouteDecision(in, d)
	return d
}
