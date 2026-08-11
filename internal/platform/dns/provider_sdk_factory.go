package dns

import (
	"context"
	"fmt"
)

// NewDefaultSDKRegistry creates the common registry. Concrete adapters are
// registered by the composition root so this package stays provider-agnostic.
func NewDefaultSDKRegistry() *SDKRegistry { return NewSDKRegistry() }

// RequireSDK returns a registered SDK and verifies its declared capabilities.
func RequireSDK(ctx context.Context, r *SDKRegistry, cfg ProviderConfig, capabilities ...string) (ProviderSDK, error) {
	if r == nil { return nil, fmt.Errorf("SDK registry is required") }
	if ctx == nil { return nil, fmt.Errorf("context is required") }
	sdk, err := r.Open(cfg)
	if err != nil { return nil, err }
	if err := ValidateProviderSDK(ctx, sdk, capabilities...); err != nil { return nil, err }
	return sdk, nil
}
