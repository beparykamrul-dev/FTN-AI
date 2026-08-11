package module

import "fmt"

type Provider struct {
	Module string `json:"module"`
	Capability string `json:"capability"`
	Priority int `json:"priority"`
	Healthy bool `json:"healthy"`
}

type ProviderSelector struct { providers []Provider }

func NewProviderSelector(providers []Provider) *ProviderSelector { return &ProviderSelector{providers: providers} }

// Select returns the highest-priority healthy provider for a capability.
// Ties are rejected to keep provider choice deterministic and explicit.
func (s *ProviderSelector) Select(capability string) (Provider, error) {
	var best Provider
	found := false
	for _, p := range s.providers {
		if p.Capability != capability || !p.Healthy { continue }
		if !found || p.Priority > best.Priority {
			best, found = p, true
		} else if p.Priority == best.Priority {
			return Provider{}, fmt.Errorf("ambiguous healthy providers for capability %q", capability)
		}
	}
	if !found { return Provider{}, fmt.Errorf("no healthy provider for capability %q", capability) }
	return best, nil
}
