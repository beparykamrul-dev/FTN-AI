package observability

// RecoveryPolicy controls when a cooled-down backend may re-enter routing.
type RecoveryPolicy struct {
	CooldownSeconds int64
	RequiredHealthyChecks uint32
}

func (p RecoveryPolicy) Valid() bool {
	return p.CooldownSeconds > 0 && p.RequiredHealthyChecks > 0
}
