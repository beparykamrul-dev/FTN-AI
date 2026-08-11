package dns

import (
	"context"
	"fmt"
)

type ProviderSDK interface {
	ProviderAdapter
	Capabilities() []string
	Health(ctx context.Context) error
	ImportZone(ctx context.Context, zone string) (Zone, error)
}

type ProviderSDKFactory func(ProviderConfig) (ProviderSDK, error)

type SDKRegistry struct {
	factories map[ProviderType]ProviderSDKFactory
}

func NewSDKRegistry() *SDKRegistry { return &SDKRegistry{factories: make(map[ProviderType]ProviderSDKFactory)} }

func (r *SDKRegistry) Register(t ProviderType, factory ProviderSDKFactory) error {
	if t == "" { return fmt.Errorf("provider type is required") }
	if factory == nil { return fmt.Errorf("provider SDK factory is required") }
	r.factories[t] = factory
	return nil
}

func (r *SDKRegistry) Open(cfg ProviderConfig) (ProviderSDK, error) {
	factory, ok := r.factories[cfg.Type]
	if !ok { return nil, fmt.Errorf("provider SDK not registered: %s", cfg.Type) }
	if !cfg.Enabled { return nil, fmt.Errorf("provider disabled: %s", cfg.ID) }
	return factory(cfg)
}

func (r *SDKRegistry) HealthAll(ctx context.Context, configs []ProviderConfig) map[string]error {
	result := make(map[string]error)
	for _, cfg := range configs {
		sdk, err := r.Open(cfg)
		if err != nil { result[cfg.ID] = err; continue }
		result[cfg.ID] = sdk.Health(ctx)
	}
	return result
}
