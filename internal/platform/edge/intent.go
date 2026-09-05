package edge

import (
	"strings"
	"time"
)

type IntentType string

const (
	IntentRoute    IntentType = "route"
	IntentFirewall IntentType = "firewall"
	IntentQoS      IntentType = "qos"
	IntentTunnel   IntentType = "tunnel"
	IntentDNS      IntentType = "dns"
)

type NetworkIntent struct {
	ID          string            `json:"id"`
	Type        IntentType        `json:"type"`
	Target      string            `json:"target"`
	Spec        map[string]string `json:"spec"`
	RequestedBy string            `json:"requested_by"`
	Approved    bool              `json:"approved"`
	CreatedAt   time.Time         `json:"created_at"`
}

func NewIntent(id string, typ IntentType, target, requestedBy string, spec map[string]string) NetworkIntent {
	copySpec := make(map[string]string, len(spec))
	for k, v := range spec {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		copySpec[k] = strings.TrimSpace(v)
	}
	return NetworkIntent{ID: strings.TrimSpace(id), Type: IntentType(strings.TrimSpace(string(typ))), Target: strings.TrimSpace(target), RequestedBy: strings.TrimSpace(requestedBy), Spec: copySpec, CreatedAt: time.Now().UTC()}
}
