package control

// ProtectionPolicy captures the minimum control-plane protections expected
// before FTN automation can affect DNS, HTTP routing, or BGP state.
type ProtectionPolicy struct {
	RequireApproval bool
	RequireAudit    bool
	RequireHealth   bool
	DryRunAllowed   bool
}

func (p ProtectionPolicy) Valid() bool {
	return p.RequireApproval && p.RequireAudit && p.RequireHealth
}
