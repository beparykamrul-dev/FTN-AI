package cctv

// SecurityPolicy applies the same FTN security boundary to routers and CCTV.
type SecurityPolicy struct { FirewallEnabled bool `json:"firewall_enabled"`; Allowlisted bool `json:"allowlisted"`; TLSRequired bool `json:"tls_required"`; AuditEnabled bool `json:"audit_enabled"`; PowerOnly bool `json:"power_only"` }

func DefaultSecurityPolicy() SecurityPolicy { return SecurityPolicy{FirewallEnabled:true, TLSRequired:true, AuditEnabled:true} }
