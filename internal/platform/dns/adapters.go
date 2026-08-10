package dns

import (
    "context"
    "fmt"
)

type AdapterFactory struct{}

func (AdapterFactory) Build(cfg ProviderConfig) (ProviderAdapter, error) {
    if !cfg.Enabled { return nil, fmt.Errorf("provider %q is disabled", cfg.ID) }
    switch cfg.Type {
    case ProviderPowerDNS, ProviderTechnitium, ProviderCoreDNS, ProviderUnbound, ProviderDNSDist,
        ProviderGoDNS, ProviderAnycast, ProviderDNSPod, ProviderCloudflare, ProviderAkamai:
        return &genericAdapter{provider: cfg.Type}, nil
    default:
        return nil, fmt.Errorf("unsupported DNS provider: %s", cfg.Type)
    }
}

type genericAdapter struct { provider ProviderType }
func (a *genericAdapter) Type() ProviderType { return a.provider }
func (a *genericAdapter) ApplyZone(context.Context, Zone) error { return nil }
func (a *genericAdapter) DeleteZone(context.Context, string) error { return nil }
