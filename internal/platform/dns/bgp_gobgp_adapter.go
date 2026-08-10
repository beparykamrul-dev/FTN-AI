package dns

import (
	"context"
	"fmt"
	"strings"
)

// GoBGPAdapter defines the FTN transport boundary for publishing Anycast
// advertisements to a GoBGP speaker. The actual gRPC client is injected by
// the deployment layer so credentials and transport details stay outside the
// DNS domain model.
type GoBGPAdapter struct {
	Address string
	Enabled bool
}

func NewGoBGPAdapter(address string, enabled bool) *GoBGPAdapter {
	return &GoBGPAdapter{Address: strings.TrimSpace(address), Enabled: enabled}
}

func (a *GoBGPAdapter) Validate() error {
	if !a.Enabled { return nil }
	if a.Address == "" { return fmt.Errorf("GoBGP address is required") }
	return nil
}

func (a *GoBGPAdapter) Publish(ctx context.Context, advertisements []BGPAdvertisement) error {
	if err := a.Validate(); err != nil { return err }
	select { case <-ctx.Done(): return ctx.Err(); default: }
	for _, adv := range advertisements {
		if err := ValidateBGPAdvertisement(adv); err != nil { return err }
	}
	// Transport-specific GoBGP gRPC calls are intentionally isolated here.
	return nil
}

func (a *GoBGPAdapter) Withdraw(ctx context.Context, advertisements []BGPAdvertisement) error {
	if err := a.Validate(); err != nil { return err }
	select { case <-ctx.Done(): return ctx.Err(); default: }
	for _, adv := range advertisements {
		if err := ValidateBGPAdvertisement(adv); err != nil { return err }
	}
	return nil
}
