package dns

import (
	"strings"
	"sync"
)

type ProviderKind string

const (
	ProviderPowerDNS ProviderKind = "powerdns"
	ProviderTechnitium ProviderKind = "technitium"
	ProviderCoreDNS ProviderKind = "coredns"
	ProviderUnbound ProviderKind = "unbound"
	ProviderDNSDist ProviderKind = "dnsdist"
	ProviderGoDNS ProviderKind = "godns"
	ProviderAnycast ProviderKind = "anycast"
	ProviderDNSPod ProviderKind = "dnspod"
	ProviderCloudflare ProviderKind = "cloudflare"
	ProviderAkamai ProviderKind = "akamai"
	ProviderFTN ProviderKind = "ftn"
)

type Provider struct {
	ID string `json:"id"`
	Kind ProviderKind `json:"kind"`
	Name string `json:"name"`
	Endpoint string `json:"endpoint,omitempty"`
	Enabled bool `json:"enabled"`
	Primary bool `json:"primary"`
	Capabilities []string `json:"capabilities,omitempty"`
}

type ProviderRegistry struct {
	mu sync.RWMutex
	providers map[string]Provider
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{providers: make(map[string]Provider)}
}

func (r *ProviderRegistry) Upsert(p Provider) bool {
	p.ID = strings.TrimSpace(p.ID)
	p.Name = strings.TrimSpace(p.Name)
	if p.ID == "" || p.Name == "" || p.Kind == "" { return false }
	r.mu.Lock(); defer r.mu.Unlock()
	r.providers[p.ID] = p
	return true
}

func (r *ProviderRegistry) Get(id string) (Provider, bool) {
	r.mu.RLock(); defer r.mu.RUnlock()
	p, ok := r.providers[strings.TrimSpace(id)]
	return p, ok
}

func (r *ProviderRegistry) List() []Provider {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers { out = append(out, p) }
	return out
}
