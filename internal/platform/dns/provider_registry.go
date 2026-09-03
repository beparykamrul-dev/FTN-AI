package dns

import (
	"sort"
	"strings"
	"sync"
)

// Provider is the provider-neutral DNS/node contract. Adding a provider or
// FTN DNS node must not require rebuilding the existing FamilyTimeNet DNS.
type Provider struct {
	ID           string       `json:"id"`
	Kind         ProviderType `json:"kind"`
	Name         string       `json:"name"`
	Endpoint     string       `json:"endpoint,omitempty"`
	Region       string       `json:"region,omitempty"`
	Scope        string       `json:"scope,omitempty"`
	Enabled      bool         `json:"enabled"`
	Primary      bool         `json:"primary"`
	Capabilities []string     `json:"capabilities,omitempty"`
	DNSZone      string       `json:"dnsZone,omitempty"`
	ConfigRef    string       `json:"configRef,omitempty"`
}

type ProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{providers: make(map[string]Provider)}
}

func (r *ProviderRegistry) Upsert(p Provider) bool {
	p.ID = strings.TrimSpace(p.ID)
	p.Name = strings.TrimSpace(p.Name)
	p.Endpoint = strings.TrimRight(strings.TrimSpace(p.Endpoint), "/")
	p.Region = strings.TrimSpace(p.Region)
	p.Scope = strings.ToLower(strings.TrimSpace(p.Scope))
	p.DNSZone = strings.TrimSpace(p.DNSZone)
	if p.ID == "" || p.Name == "" || p.Kind == "" || p.DNSZone == "" {
		return false
	}
	if p.Scope != "local" && p.Scope != "global" {
		return false
	}
	r.mu.Lock()
	r.providers[p.ID] = p
	r.mu.Unlock()
	return true
}

func (r *ProviderRegistry) Get(id string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[strings.TrimSpace(id)]
	return p, ok
}

func (r *ProviderRegistry) List() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *ProviderRegistry) Select(scope, zone string) []Provider {
	scope = strings.ToLower(strings.TrimSpace(scope))
	zone = strings.TrimSpace(zone)
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Provider, 0)
	for _, p := range r.providers {
		if p.Enabled && p.Scope == scope && p.DNSZone == zone {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Primary != out[j].Primary {
			return out[i].Primary
		}
		return out[i].ID < out[j].ID
	})
	return out
}
