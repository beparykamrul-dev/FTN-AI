package edge

import("context";"strings")
type FirewallRule struct { ID string `json:"id"`; Family string `json:"family"`; Chain string `json:"chain"`; Action string `json:"action"`; Source string `json:"source,omitempty"`; Destination string `json:"destination,omitempty"`; Comment string `json:"comment,omitempty"` }
type LocalRouter interface { CoreRouter; FirewallRules(context.Context)([]FirewallRule,error); ApplyFirewall(context.Context,FirewallRule)error }
func (r FirewallRule) Valid() bool { return strings.TrimSpace(r.ID)!=""&&len(strings.TrimSpace(r.ID))<=256&&strings.TrimSpace(r.Family)!=""&&len(strings.TrimSpace(r.Family))<=32&&strings.TrimSpace(r.Chain)!=""&&len(strings.TrimSpace(r.Chain))<=64&&strings.TrimSpace(r.Action)!=""&&len(strings.TrimSpace(r.Action))<=64&&len(strings.TrimSpace(r.Source))<=256&&len(strings.TrimSpace(r.Destination))<=256&&len(strings.TrimSpace(r.Comment))<=1024 }
