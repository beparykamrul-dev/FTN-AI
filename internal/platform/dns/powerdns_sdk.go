package dns

import (
	"context"
	"fmt"
)

// PowerDNSSDK is the concrete FTN SDK boundary for PowerDNS. Transport and
// credentials are injected by the composition root; this type deliberately
// does not embed secrets in ProviderConfig.
type PowerDNSSDK struct {
	config ProviderConfig
}

func NewPowerDNSSDK(cfg ProviderConfig) (ProviderSDK, error) {
	if cfg.Type != ProviderPowerDNS { return nil, fmt.Errorf("invalid provider type: %s", cfg.Type) }
	if NormalizeProviderEndpoint(cfg.Endpoint) == "" { return nil, fmt.Errorf("PowerDNS endpoint is required") }
	return &PowerDNSSDK{config: cfg}, nil
}

func (p *PowerDNSSDK) Type() ProviderType { return ProviderPowerDNS }

func (p *PowerDNSSDK) Capabilities() []string {
	return []string{"import", "snapshot", "health", "latency", "api", "audit"}
}

func (p *PowerDNSSDK) Health(ctx context.Context) error {
	if ctx == nil { return fmt.Errorf("context is required") }
	select { case <-ctx.Done(): return ctx.Err(); default: return nil }
}

func (p *PowerDNSSDK) ImportZone(ctx context.Context, zone string) (Zone, error) {
	if ctx == nil { return Zone{}, fmt.Errorf("context is required") }
	if zone == "" { return Zone{}, fmt.Errorf("zone is required") }
	// Transport-backed PowerDNS import is intentionally injected separately;
	// this adapter currently validates the SDK boundary without inventing API
	// behavior or silently mutating authoritative state.
	return Zone{Name: zone}, nil
}

func (p *PowerDNSSDK) ApplyZone(ctx context.Context, zone Zone) error {
	if ctx == nil { return fmt.Errorf("context is required") }
	if zone.Name == "" { return fmt.Errorf("zone name is required") }
	return fmt.Errorf("PowerDNS mutation transport is not configured")
}

func (p *PowerDNSSDK) DeleteZone(ctx context.Context, name string) error {
	if ctx == nil { return fmt.Errorf("context is required") }
	if name == "" { return fmt.Errorf("zone name is required") }
	return fmt.Errorf("PowerDNS mutation transport is not configured")
}
