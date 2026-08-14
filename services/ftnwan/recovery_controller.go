package ftnwan

// RecoveryPlan is a provider-neutral, non-executing recovery plan.
type RecoveryPlan struct {
	NodeID        string
	TargetVersion string
	Reason        string
	Approved      bool
}

// BuildRecoveryPlan creates a safe plan only when the recovery policy allows it.
// Execution belongs to the transaction/reconciliation controller.
func BuildRecoveryPlan(policy RecoveryPolicy, nodeID, targetVersion, reason string, hasKnownGoodState bool) (RecoveryPlan, RecoveryDecision) {
	if EvaluateRecovery(policy, hasKnownGoodState) != RecoveryAllow || nodeID == "" || targetVersion == "" {
		return RecoveryPlan{}, RecoveryDeny
	}

	return RecoveryPlan{
		NodeID:        nodeID,
		TargetVersion: targetVersion,
		Reason:        reason,
		Approved:      !policy.RequireApproval,
	}, RecoveryAllow
}
