package dns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type BGPHealthState struct {
	NodeID string `json:"node_id"`
	Prefix string `json:"prefix"`
	Healthy bool `json:"healthy"`
	LastCheck time.Time `json:"last_check"`
}

// BGPHealthWithdrawal coordinates health state with advertisement withdrawal.
// The caller supplies the BGP transport, keeping routing credentials outside
// the DNS domain model.
type BGPHealthWithdrawal struct {
	mu sync.RWMutex
	states map[string]BGPHealthState
}

func NewBGPHealthWithdrawal() *BGPHealthWithdrawal {
	return &BGPHealthWithdrawal{states: make(map[string]BGPHealthState)}
}

func (c *BGPHealthWithdrawal) Update(state BGPHealthState) {
	c.mu.Lock(); defer c.mu.Unlock()
	c.states[state.NodeID+"|"+state.Prefix] = state
}

func (c *BGPHealthWithdrawal) Healthy(nodeID, prefix string) bool {
	c.mu.RLock(); defer c.mu.RUnlock()
	s, ok := c.states[nodeID+"|"+prefix]
	return ok && s.Healthy
}

func (c *BGPHealthWithdrawal) Reconcile(ctx context.Context, adapter *GoBGPAdapter, advertisements []BGPAdvertisement) error {
	if adapter == nil { return fmt.Errorf("GoBGP adapter is required") }
	if err := adapter.Validate(); err != nil { return err }
	var withdraw []BGPAdvertisement
	for _, adv := range advertisements {
		if !c.Healthy(adv.NodeID, adv.Prefix) { withdraw = append(withdraw, adv) }
	}
	if len(withdraw) == 0 { return nil }
	return adapter.Withdraw(ctx, withdraw)
}
