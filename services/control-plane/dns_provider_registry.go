package main

import (
	"sort"
	"strings"
)

type DNSProviderEndpoint struct {
	Protocol   string `json:"protocol" yaml:"protocol"`
	Address    string `json:"address" yaml:"address"`
	ServerName string `json:"server_name,omitempty" yaml:"server_name,omitempty"`
}

type DNSProvider struct {
	ID        string               `json:"id" yaml:"id"`
	Name      string               `json:"name" yaml:"name"`
	Enabled   bool                 `json:"enabled" yaml:"enabled"`
	Endpoints []DNSProviderEndpoint `json:"endpoints" yaml:"endpoints"`
}

type DNSProviderHealth struct {
	ProviderID string  `json:"provider_id"`
	Endpoint   string  `json:"endpoint"`
	Protocol   string  `json:"protocol"`
	Healthy    bool    `json:"healthy"`
	LatencyMS  float64 `json:"latency_ms"`
	Failures   int     `json:"failures"`
}

// RankDNSProviders returns healthy providers in ascending observed latency.
// The resolver can use the first N entries for parallel racing/failover.
func RankDNSProviders(health []DNSProviderHealth, maxLatencyMS float64) []DNSProviderHealth {
	out := make([]DNSProviderHealth, 0, len(health))
	for _, h := range health {
		if !h.Healthy || h.LatencyMS < 0 {
			continue
		}
		if maxLatencyMS > 0 && h.LatencyMS > maxLatencyMS {
			continue
		}
		out = append(out, h)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].LatencyMS == out[j].LatencyMS {
			return out[i].ProviderID+"/"+out[i].Endpoint < out[j].ProviderID+"/"+out[j].Endpoint
		}
		return out[i].LatencyMS < out[j].LatencyMS
	})
	return out
}

// PreferredDNSEndpoints filters a provider's endpoints by the configured
// protocol preference while preserving the preference order.
func PreferredDNSEndpoints(provider DNSProvider, protocols []string) []DNSProviderEndpoint {
	if !provider.Enabled {
		return nil
	}
	if len(protocols) == 0 {
		return append([]DNSProviderEndpoint(nil), provider.Endpoints...)
	}
	result := make([]DNSProviderEndpoint, 0, len(provider.Endpoints))
	for _, protocol := range protocols {
		want := strings.ToLower(strings.TrimSpace(protocol))
		for _, endpoint := range provider.Endpoints {
			if strings.ToLower(endpoint.Protocol) == want {
				result = append(result, endpoint)
			}
		}
	}
	return result
}
