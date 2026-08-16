package security

// ApprovalExpiry describes the validity window of an explicit deployment approval.
type ApprovalExpiry struct {
	ApprovedBy string
	ApprovedUnix int64
	ExpiresUnix int64
}

func (a ApprovalExpiry) Valid(nowUnix int64) bool {
	return a.ApprovedBy != "" && a.ExpiresUnix > a.ApprovedUnix && nowUnix < a.ExpiresUnix
}
