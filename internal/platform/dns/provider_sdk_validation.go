package dns

import (
	"context"
	"fmt"
)

// ValidateProviderSDK performs the common runtime checks before a concrete
// provider SDK is admitted to a control-plane operation.
func ValidateProviderSDK(ctx context.Context, sdk ProviderSDK, requiredCapabilities ...string) error {
	if sdk == nil { return fmt.Errorf("provider SDK is required") }
	select { case <-ctx.Done(): return ctx.Err(); default: }

	caps := make(map[string]struct{}, len(sdk.Capabilities()))
	for _, capability := range sdk.Capabilities() { caps[capability] = struct{}{} }
	for _, required := range requiredCapabilities {
		if _, ok := caps[required]; !ok { return fmt.Errorf("provider %s does not support capability %q", sdk.Type(), required) }
	}
	return nil
}
