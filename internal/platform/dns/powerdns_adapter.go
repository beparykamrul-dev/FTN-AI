package dns

import (
	"context"
	"fmt"
)

// PowerDNSAdapter is the FTN integration boundary for PowerDNS Authoritative.
// Network I/O and credentials are injected by the application layer; this
// package only defines the validated zone operation contract.
type PowerDNSAdapter struct {
	Endpoint string
}

func NewPowerDNSAdapter(endpoint string) *PowerDNSAdapter {
	return &PowerDNSAdapter{Endpoint: NormalizeProviderEndpoint(endpoint)}
}

func (a *PowerDNSAdapter) Type() ProviderType { return ProviderPowerDNS }

func (a *PowerDNSAdapter) ApplyZone(ctx context.Context, zone Zone) error {
	if a == nil || a.Endpoint == "" {
		return fmt.Errorf("PowerDNS endpoint is required")
	}
	if zone.Name == "" {
		return fmt.Errorf("zone name is required")
	}
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (a *PowerDNSAdapter) DeleteZone(ctx context.Context, zone string) error {
	if a == nil || a.Endpoint == "" {
		return fmt.Errorf("PowerDNS endpoint is required")
	}
	if normalizeZone(zone) == "" {
		return fmt.Errorf("zone name is required")
	}
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}
