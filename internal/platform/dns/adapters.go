package dns

import (
    "context"
    "fmt"
)

type AdapterFactory struct{}

func (AdapterFactory) Build(cfg ProviderConfig) (ProviderAdapter, error) {
    if !cfg.Enabled { return nil, fmt.Errorf("provider %q is disabled", cfg.ID) }
    if len(cfg.ID) == 0 || len(cfg.ID) > 128 { return nil, fmt.Errorf("invalid provider ID") }
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
func (a *genericAdapter) ApplyZone(ctx context.Context, _ Zone) error {
    if ctx == nil { return fmt.Errorf("context is required") }
    return fmt.Errorf("DNS provider %q has no mutation implementation", a.provider)
}
func (a *genericAdapter) DeleteZone(ctx context.Context, _ string) error {
    if ctx == nil { return fmt.Errorf("context is required") }
    return fmt.Errorf("DNS provider %q has no mutation implementation", a.provider)
}
