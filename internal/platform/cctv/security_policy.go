package cctv

type SecurityPolicy struct{FirewallEnabled bool `json:"firewall_enabled"`;Allowlisted bool `json:"allowlisted"`;TLSRequired bool `json:"tls_required"`;AuditEnabled bool `json:"audit_enabled"`;PowerOnly bool `json:"power_only"`}
func(p SecurityPolicy)Valid()bool{return p.FirewallEnabled&&p.Allowlisted&&p.TLSRequired&&p.AuditEnabled}
func DefaultSecurityPolicy()SecurityPolicy{return SecurityPolicy{FirewallEnabled:true,Allowlisted:true,TLSRequired:true,AuditEnabled:true}}
