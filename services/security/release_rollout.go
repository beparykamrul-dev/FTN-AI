package security

// RolloutDecision is the security result consumed by deployment orchestration.
type RolloutDecision struct {
	ReleaseID string
	Allowed bool
	Reason string
}

func EvaluateRollout(v ReleaseVerification, policy ReleasePolicy, score int) RolloutDecision {
	if !v.Ready() { return RolloutDecision{ReleaseID: v.ReleaseID, Allowed: false, Reason: "release-verification-failed"} }
	if !policy.Allows(v, score) { return RolloutDecision{ReleaseID: v.ReleaseID, Allowed: false, Reason: "release-policy-failed"} }
	return RolloutDecision{ReleaseID: v.ReleaseID, Allowed: true, Reason: "release-security-approved"}
}
