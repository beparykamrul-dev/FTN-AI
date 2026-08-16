package security

// ApprovalRecord captures an explicit security approval decision.
type ApprovalRecord struct {
	ChangeID    string
	ApprovedBy  string
	Reason      string
	ExpiresUnix int64
}

func (a ApprovalRecord) Valid(nowUnix int64) bool {
	return a.ChangeID != "" && a.ApprovedBy != "" && a.Reason != "" && a.ExpiresUnix > nowUnix
}
