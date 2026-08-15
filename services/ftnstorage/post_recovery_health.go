package ftnstorage

// PostRecoveryHealth is the final health verdict after repair or rollback.
type PostRecoveryHealth struct {
	NodeID   string
	Verified bool
	Healthy  bool
	Rejoin   bool
	Reason   string
}

func EvaluatePostRecovery(nodeID string, verified bool) PostRecoveryHealth {
	return PostRecoveryHealth{
		NodeID:   nodeID,
		Verified: verified,
		Healthy:  verified,
		Rejoin:   verified,
		Reason:   "post-recovery-integrity-and-health-check",
	}
}
