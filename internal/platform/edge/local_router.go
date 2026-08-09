package edge

import "context"

type FirewallRule struct { ID string `json:"id"`; Family string `json:"family"`; Chain string `json:"chain"`; Action string `json:"action"`; Source string `json:"source,omitempty"`; Destination string `json:"destination,omitempty"`; Comment string `json:"comment,omitempty"` }

type LocalRouter interface {
	CoreRouter
	FirewallRules(ctx context.Context) ([]FirewallRule, error)
	ApplyFirewall(ctx context.Context, rule FirewallRule) error
}

// The concrete Linux implementation should use nftables/netlink through a
// privileged FTN agent. The control plane only carries approved intent.
