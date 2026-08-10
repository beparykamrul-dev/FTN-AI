package dns

import (
	"context"
	"fmt"
	"strings"
)

type ProviderType string

const (
	ProviderPowerDNS ProviderType = "powerdns"
	ProviderTechnitium ProviderType = "technitium"
	ProviderCoreDNS ProviderType = "coredns"
	ProviderUnbound ProviderType = "unbound"
	ProviderDNSDist ProviderType = "dnsdist"
	ProviderGoDNS ProviderType = "godns"
	ProviderAnycast ProviderType = "anycast"
	ProviderDNSPod ProviderType = "dnspod"
	ProviderCloudflare ProviderType = "cloudflare"
	ProviderAkamai ProviderType = "akamai"
)

type ProviderConfig struct {
	ID string `json:"id"`
	Type ProviderType `json:"type"`
	Endpoint string `json:"endpoint,omitempty"`
	Enabled bool `json:"enabled"`
}

type ProviderAdapter interface {
	Type() ProviderType
	ApplyZone(context.Context, Zone) error
	DeleteZone(context.Context, string) error
}

// AdapterRegistry keeps provider-specific DNS implementations behind one
// control-plane interface. Credentials are intentionally excluded from config.
type AdapterRegistry struct { adapters map[ProviderType]ProviderAdapter }

func NewAdapterRegistry() *AdapterRegistry { return &AdapterRegistry{adapters: make(map[ProviderType]ProviderAdapter)} }

func (r *AdapterRegistry) Register(adapter ProviderAdapter) error {
	if adapter == nil { return fmt.Errorf("adapter is required") }
	r.adapters[adapter.Type()] = adapter
	return nil
}

func (r *AdapterRegistry) Get(t ProviderType) (ProviderAdapter, bool) { a, ok := r.adapters[t]; return a, ok }

func NormalizeProviderEndpoint(endpoint string) string {
	return strings.TrimRight(strings.TrimSpace(endpoint), "/")
}
