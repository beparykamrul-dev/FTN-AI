package dns

import (
	"context"
	"fmt"
	"strings"
)

// ProviderType is retained as a compatibility alias for the canonical
// ProviderKind registry model.
type ProviderType = ProviderKind

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
