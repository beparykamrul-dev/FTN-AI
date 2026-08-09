package controlplane

import "time"

type PolicyScope string
const (
	ScopeDNS PolicyScope = "dns"
	ScopeNetwork PolicyScope = "network"
	ScopeHTTP PolicyScope = "http"
	ScopeResolver PolicyScope = "resolver"
	ScopeInfrastructure PolicyScope = "infrastructure"
)

type Policy struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Scope PolicyScope `json:"scope"`
	TenantID string `json:"tenant_id"`
	Priority int `json:"priority"`
	Enabled bool `json:"enabled"`
	Rules []Rule `json:"rules,omitempty"`
	Version uint64 `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rule struct {
	Action string `json:"action"`
	Match map[string]string `json:"match,omitempty"`
}

type PolicyDecision struct {
	PolicyID string `json:"policy_id"`
	Action string `json:"action"`
	Reason string `json:"reason"`
}

// Evaluate is deterministic policy evaluation. Runtime enforcement belongs to
// the authorized data-plane adapter; the control plane records the decision.
func Evaluate(p Policy, attrs map[string]string) PolicyDecision {
	if !p.Enabled { return PolicyDecision{PolicyID:p.ID, Action:"allow", Reason:"policy_disabled"} }
	for _, r := range p.Rules {
		matched := true
		for k, v := range r.Match { if attrs[k] != v { matched = false; break } }
		if matched { return PolicyDecision{PolicyID:p.ID, Action:r.Action, Reason:"rule_match"} }
	}
	return PolicyDecision{PolicyID:p.ID, Action:"allow", Reason:"no_rule_match"}
}
