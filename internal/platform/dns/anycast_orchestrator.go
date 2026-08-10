package dns

import (
	"context"
	"fmt"
)

type AnycastOrchestrator struct {
	Policy *ProbePolicyController
	Health *BGPHealthWithdrawal
	BGP *GoBGPAdapter
}

func NewAnycastOrchestrator(policy *ProbePolicyController, health *BGPHealthWithdrawal, bgp *GoBGPAdapter) *AnycastOrchestrator {
	return &AnycastOrchestrator{Policy: policy, Health: health, BGP: bgp}
}

// Reconcile consumes real DNS probe results and produces the BGP state change
// required by the current health policy. It does not mutate router state when
// dependencies are missing.
func (o *AnycastOrchestrator) Reconcile(ctx context.Context, key string, adv BGPAdvertisement, probe DNSProbeResult, policy ProbePolicy) error {
	if o.Policy == nil || o.Health == nil || o.BGP == nil { return fmt.Errorf("anycast orchestrator dependencies are required") }
	healthy, err := o.Policy.Observe(ctx, key, probe, policy)
	if err != nil { return err }
	o.Health.Update(BGPHealthState{NodeID: adv.NodeID, Prefix: adv.Prefix, Healthy: healthy})
	if healthy { return o.BGP.Publish(ctx, []BGPAdvertisement{adv}) }
	return o.BGP.Withdraw(ctx, []BGPAdvertisement{adv})
}
