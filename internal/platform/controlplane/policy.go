package controlplane

import("fmt";"strings";"time")
type PolicyScope string
const(ScopeDNS PolicyScope="dns";ScopeNetwork PolicyScope="network";ScopeHTTP PolicyScope="http";ScopeResolver PolicyScope="resolver";ScopeInfrastructure PolicyScope="infrastructure")
type Policy struct{ID string `json:"id"`;Name string `json:"name"`;Scope PolicyScope `json:"scope"`;TenantID string `json:"tenant_id"`;Priority int `json:"priority"`;Enabled bool `json:"enabled"`;Rules []Rule `json:"rules,omitempty"`;Version uint64 `json:"version"`;UpdatedAt time.Time `json:"updated_at"`}
type Rule struct{Action string `json:"action"`;Match map[string]string `json:"match,omitempty"`}
type PolicyDecision struct{PolicyID string `json:"policy_id"`;Action string `json:"action"`;Reason string `json:"reason"`}
func(p Policy)Validate()error{p.ID=strings.TrimSpace(p.ID);p.TenantID=strings.TrimSpace(p.TenantID);if p.ID==""||p.TenantID==""{return fmt.Errorf("policy id and tenant id are required")};switch p.Scope{case ScopeDNS,ScopeNetwork,ScopeHTTP,ScopeResolver,ScopeInfrastructure:default:return fmt.Errorf("invalid policy scope")};if len(p.Rules)>1024{return fmt.Errorf("too many policy rules")};for _,r:=range p.Rules{if strings.TrimSpace(r.Action)==""||len(r.Action)>64{return fmt.Errorf("invalid policy action")};if len(r.Match)>256{return fmt.Errorf("too many policy match attributes")}};return nil}
func Evaluate(p Policy,attrs map[string]string)PolicyDecision{if !p.Enabled{return PolicyDecision{PolicyID:strings.TrimSpace(p.ID),Action:"allow",Reason:"policy_disabled"}};for _,r:=range p.Rules{matched:=true;for k,v:=range r.Match{if attrs==nil||attrs[k]!=v{matched=false;break}};if matched{return PolicyDecision{PolicyID:strings.TrimSpace(p.ID),Action:strings.ToLower(strings.TrimSpace(r.Action)),Reason:"rule_match"}}};return PolicyDecision{PolicyID:strings.TrimSpace(p.ID),Action:"allow",Reason:"no_rule_match"}}
