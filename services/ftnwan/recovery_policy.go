package ftnwan

// RecoveryPolicy defines the provider-neutral recovery boundary for FTNWAN.
// It describes what may be recovered; it does not perform recovery itself.
type RecoveryPolicy struct {
	Name               string
	StateRetentionDays int
	RequireKnownGood   bool
	AllowAutoFailover  bool
	RequireApproval    bool
}

// RecoveryDecision is the deterministic baseline decision returned by the
// recovery policy layer before a deployment/recovery controller acts.
type RecoveryDecision string

const (
	RecoveryAllow RecoveryDecision = "allow"
	RecoveryDeny  RecoveryDecision = "deny"
)

// EvaluateRecovery checks only safety prerequisites. Actual restore operations
// remain in the transaction/reconciliation layer.
func EvaluateRecovery(p RecoveryPolicy, hasKnownGoodState bool) RecoveryDecision {
	if p.Name == "" || p.StateRetentionDays <= 0 {
		return RecoveryDeny
	}
	if p.RequireKnownGood && !hasKnownGoodState {
		return RecoveryDeny
	}
	return RecoveryAllow
}
