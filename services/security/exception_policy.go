package security

// SecurityException represents an explicit, auditable exception to a finding.
type SecurityException struct {
	RuleID      string
	Reason      string
	ApprovedBy  string
	ExpiresUnix int64
}

func (e SecurityException) Valid(nowUnix int64) bool {
	return e.RuleID != "" && e.Reason != "" && e.ApprovedBy != "" && e.ExpiresUnix > nowUnix
}
