package main

import (
	"errors"
	"math"
	"sort"
	"strings"
)

type TunnelProtocol struct {
	ID       string   `json:"id"`
	Family   string   `json:"family"`
	Transports []string `json:"transports"`
	Enabled  bool     `json:"enabled"`
}

type TunnelEndpoint struct {
	ID        string `json:"id"`
	Protocol  string `json:"protocol"`
	Transport string `json:"transport"`
	Address   string `json:"address"`
	Healthy   bool   `json:"healthy"`
	LatencyMS float64 `json:"latency_ms"`
	MTU       int    `json:"mtu,omitempty"`
}

type MultiTunnelPolicy struct {
	Protocols       []TunnelProtocol `json:"protocols"`
	MaxActivePaths  int              `json:"max_active_paths"`
	Failover        bool             `json:"failover"`
	HealthRequired  bool             `json:"health_required"`
	FTNOwnedOnly    bool             `json:"ftn_owned_only"`
}

func ValidateMultiTunnelPolicy(p MultiTunnelPolicy) error {
	if p.MaxActivePaths < 1 || p.MaxActivePaths > 32 { return errors.New("max_active_paths_out_of_range") }
	if !p.FTNOwnedOnly { return errors.New("ftn_owned_only_required") }
	seen := map[string]bool{}
	for _, proto := range p.Protocols {
		id := strings.ToLower(strings.TrimSpace(proto.ID))
		if id == "" || seen[id] { return errors.New("invalid_or_duplicate_protocol") }
		seen[id] = true
	}
	return nil
}

// SelectMultiTunnelPaths returns healthy FTN-approved paths ordered for
// deterministic active/standby use. It never attempts to circumvent an
// external restriction or modify a network policy.
func SelectMultiTunnelPaths(p MultiTunnelPolicy, endpoints []TunnelEndpoint) []TunnelEndpoint {
	if p.MaxActivePaths <= 0 { return nil }
	allowed := map[string]bool{}
	for _, proto := range p.Protocols { if proto.Enabled { allowed[strings.ToLower(strings.TrimSpace(proto.ID))] = true } }
	out := make([]TunnelEndpoint, 0, len(endpoints))
	for _, e := range endpoints {
		if !e.Healthy || !allowed[strings.ToLower(strings.TrimSpace(e.Protocol))] { continue }
		if math.IsNaN(e.LatencyMS) || math.IsInf(e.LatencyMS, 0) || e.LatencyMS < 0 { continue }
		out = append(out, e)
	}
	sort.SliceStable(out, func(i,j int) bool {
		if out[i].LatencyMS != out[j].LatencyMS { return out[i].LatencyMS < out[j].LatencyMS }
		return strings.TrimSpace(out[i].ID) < strings.TrimSpace(out[j].ID)
	})
	if len(out) > p.MaxActivePaths { out = out[:p.MaxActivePaths] }
	return out
}
