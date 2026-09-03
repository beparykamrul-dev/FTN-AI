package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/netip"
	"sort"
	"strings"
)

type FTNRoute struct {
	Prefix string `json:"prefix"`
	NextHop string `json:"next_hop,omitempty"`
	Protocol string `json:"protocol"`
	VRF string `json:"vrf,omitempty"`
	LocalPref uint32 `json:"local_pref,omitempty"`
	MED uint32 `json:"med,omitempty"`
	Metric uint32 `json:"metric,omitempty"`
	ASPath []uint32 `json:"as_path,omitempty"`
	Community []string `json:"community,omitempty"`
	Active bool `json:"active"`
}

type FTNRouteIntent struct {
	Action string `json:"action"`
	Route FTNRoute `json:"route"`
	Reason string `json:"reason,omitempty"`
}

type FTNRouteDecision struct {
	Allowed bool `json:"allowed"`
	RequiresApproval bool `json:"requires_approval"`
	Risk string `json:"risk"`
	Reason string `json:"reason"`
	DecisionHash string `json:"decision_hash"`
}

func NormalizeFTNRoute(r FTNRoute) (FTNRoute, error) {
	r.Prefix = strings.TrimSpace(r.Prefix)
	p, err := netip.ParsePrefix(r.Prefix)
	if err != nil { return FTNRoute{}, fmt.Errorf("invalid prefix: %w", err) }
	r.Prefix = p.Masked().String()
	r.NextHop = strings.TrimSpace(r.NextHop)
	if r.NextHop != "" { if _, err := netip.ParseAddr(r.NextHop); err != nil { return FTNRoute{}, fmt.Errorf("invalid next-hop: %w", err) } }
	r.Protocol = strings.ToLower(strings.TrimSpace(r.Protocol))
	switch r.Protocol { case "bgp", "ibgp", "ebgp", "ospf", "static", "connected": default: return FTNRoute{}, fmt.Errorf("unsupported routing protocol") }
	if r.VRF == "" { r.VRF = "default" }
	sort.Strings(r.Community)
	return r, nil
}

func HashFTNRouteIntent(in FTNRouteIntent) string {
	r, err := NormalizeFTNRoute(in.Route)
	if err != nil { r = FTNRoute{} }
	v := fmt.Sprintf("%s|%s|%s|%s|%d|%d|%d|%v|%v|%t", strings.ToLower(strings.TrimSpace(in.Action)), r.Prefix, r.NextHop, r.Protocol, r.LocalPref, r.MED, r.Metric, r.ASPath, r.Community, r.Active)
	s := sha256.Sum256([]byte(v))
	return hex.EncodeToString(s[:])
}

func EvaluateFTNRouteIntent(in FTNRouteIntent) FTNRouteDecision {
	r, err := NormalizeFTNRoute(in.Route)
	if err != nil { return FTNRouteDecision{Reason: err.Error(), Risk: "invalid", DecisionHash: HashFTNRouteIntent(in)} }
	action := strings.ToLower(strings.TrimSpace(in.Action))
	if action == "" { return FTNRouteDecision{Reason: "action_required", Risk: "invalid", DecisionHash: HashFTNRouteIntent(in)} }
	if action == "advertise-external" || action == "withdraw-external" { return FTNRouteDecision{RequiresApproval:true, Risk:"high", Reason:"external_route_change_requires_approval", DecisionHash:HashFTNRouteIntent(in)} }
	if action == "delete" || action == "flush" { return FTNRouteDecision{RequiresApproval:true, Risk:"critical", Reason:"destructive_route_change", DecisionHash:HashFTNRouteIntent(in)} }
	if !r.Active && action == "install" { return FTNRouteDecision{RequiresApproval:true, Risk:"medium", Reason:"inactive_route_requires_review", DecisionHash:HashFTNRouteIntent(in)} }
	return FTNRouteDecision{Allowed:true, RequiresApproval:true, Risk:"medium", Reason:"route_intent_validated; execution remains approval-gated", DecisionHash:HashFTNRouteIntent(in)}
}
