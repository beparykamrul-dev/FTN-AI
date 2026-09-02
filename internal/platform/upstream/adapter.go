package upstream

import (
	"context"
	"errors"
	"sort"
	"time"
)

// Adapter is the control-plane boundary for an FTN upstream connection.
// Implementations may integrate with BGP, provider APIs, IX fabrics, or
// approved CDN/interconnect programs. Credentials and secrets stay outside
// this package.
type Adapter interface {
	Name() string
	Validate(context.Context, Profile) error
	Discover(context.Context) ([]Endpoint, error)
	Health(context.Context, Endpoint) (Health, error)
}

type Profile struct {
	Name             string            `json:"name" yaml:"name"`
	Enabled          bool              `json:"enabled" yaml:"enabled"`
	ASN              uint32            `json:"asn" yaml:"asn"`
	LocalAddress     string            `json:"local_address" yaml:"local_address"`
	RemoteAddresses  []string          `json:"remote_addresses" yaml:"remote_addresses"`
	RemoteASNs       []uint32          `json:"remote_asns" yaml:"remote_asns"`
	Role             string            `json:"role" yaml:"role"`
	Communities      []string          `json:"communities" yaml:"communities"`
	AllowedPrefixes  []string          `json:"allowed_prefixes" yaml:"allowed_prefixes"`
	MaxPrefixCount   uint32            `json:"max_prefix_count" yaml:"max_prefix_count"`
	HealthInterval   time.Duration     `json:"health_interval" yaml:"health_interval"`
	ProviderMetadata map[string]string `json:"provider_metadata" yaml:"provider_metadata"`
}

type Endpoint struct {
	Name      string `json:"name" yaml:"name"`
	Address   string `json:"address" yaml:"address"`
	ASN       uint32 `json:"asn" yaml:"asn"`
	Transport string `json:"transport" yaml:"transport"`
	Purpose   string `json:"purpose" yaml:"purpose"`
}

type Health struct {
	Healthy       bool          `json:"healthy"`
	Latency       time.Duration `json:"latency"`
	PacketLossPct float64       `json:"packet_loss_pct"`
	CheckedAt     time.Time     `json:"checked_at"`
	Message       string        `json:"message,omitempty"`
}

func ValidateProfile(p Profile) error {
	if p.Name == "" {
		return errors.New("upstream profile name is required")
	}
	if !p.Enabled {
		return nil
	}
	if p.ASN == 0 {
		return errors.New("enabled upstream profile requires a local ASN")
	}
	if len(p.RemoteAddresses) != len(p.RemoteASNs) {
		return errors.New("remote_addresses and remote_asns must have equal lengths")
	}
	if p.MaxPrefixCount == 0 {
		return errors.New("max_prefix_count must be set for an enabled upstream")
	}
	return nil
}

// RankHealthyEndpoints returns deterministic candidates for failover. Lower
// latency is preferred, then lower packet loss, then name.
func RankHealthyEndpoints(endpoints []Endpoint, health map[string]Health) []Endpoint {
	out := make([]Endpoint, 0, len(endpoints))
	for _, ep := range endpoints {
		if h, ok := health[ep.Name]; ok && h.Healthy {
			out = append(out, ep)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		hi, hj := health[out[i].Name], health[out[j].Name]
		if hi.Latency != hj.Latency {
			return hi.Latency < hj.Latency
		}
		if hi.PacketLossPct != hj.PacketLossPct {
			return hi.PacketLossPct < hj.PacketLossPct
		}
		return out[i].Name < out[j].Name
	})
	return out
}
