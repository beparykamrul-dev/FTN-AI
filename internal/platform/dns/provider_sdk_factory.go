package dns

import "fmt"

// NewDefaultSDKRegistry creates the common registry. Concrete adapters are
// registered by the composition root so this package stays provider-agnostic.
func NewDefaultSDKRegistry() *SDKRegistry { return NewSDKRegistry() }

// RequireSDK returns a registered SDK and verifies its declared capabilities.
func RequireSDK(r *SDKRegistry, cfg ProviderConfig, capabilities ...string) (ProviderSDK, error) {
	if r == nil { return nil, fmt.Errorf("SDK registry is required") }
	sdk, err := r.Open(cfg)
	if err != nil { return nil, err }
	if err := ValidateProviderSDK(nilContext{}, sdk, capabilities...); err != nil { return nil, err }
	return sdk, nil
}

type nilContext struct{}
func (nilContext) Deadline() (deadlineTime, bool) { return deadlineTime{}, false }
func (nilContext) Done() <-chan struct{} { return nil }
func (nilContext) Err() error { return nil }
func (nilContext) Value(any) any { return nil }

type deadlineTime struct{}
