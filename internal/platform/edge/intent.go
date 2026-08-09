package edge

import "time"

type IntentType string
const (
	IntentRoute IntentType = "route"
	IntentFirewall IntentType = "firewall"
	IntentQoS IntentType = "qos"
	IntentTunnel IntentType = "tunnel"
	IntentDNS IntentType = "dns"
)

type NetworkIntent struct {
	ID string `json:"id"`
	Type IntentType `json:"type"`
	Target string `json:"target"`
	Spec map[string]string `json:"spec"`
	RequestedBy string `json:"requested_by"`
	Approved bool `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
}

// Intent is a provider-neutral desired state. Executors remain separate so
// RouterOS, Linux/netlink, nftables, WireGuard and vendor OLT drivers can be
// swapped without changing the control-plane model.
func NewIntent(id string, typ IntentType, target, requestedBy string, spec map[string]string) NetworkIntent {
	return NetworkIntent{ID:id, Type:typ, Target:target, RequestedBy:requestedBy, Spec:spec, CreatedAt:time.Now().UTC()}
}
