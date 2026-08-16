package security

// RiskPolicy defines deployment thresholds for aggregated security risk.
type RiskPolicy struct {
	MaxCritical uint32
	MaxHigh uint32
	RequireApproval bool
}

func (p RiskPolicy) Allows(r RiskAggregate, approved bool) bool {
	if r.Critical > p.MaxCritical || r.High > p.MaxHigh { return false }
	if p.RequireApproval && !approved { return false }
	return true
}
