package fiber

import "context"

type RecoveryDecision string
const (
	DecisionObserve RecoveryDecision = "observe"
	DecisionAutoMitigate RecoveryDecision = "auto-mitigate"
	DecisionApproval RecoveryDecision = "approval-required"
	DecisionFieldRepair RecoveryDecision = "field-repair"
)

type RecoveryPolicy struct {
	MinorDamageMaxPercent float64
	MaxLossDb float64
	MinConfidence float64
	RequireVerification bool
	AllowAutomaticReroute bool
}

type RecoveryAssessment struct {
	CoreID string
	DamagePercent float64
	Decision RecoveryDecision
	Actions []RecoveryAction
	Confidence float64
	Reason string
}

// PolicyEngine separates AI recommendations from actual privileged recovery.
type PolicyEngine interface {
	Assess(context.Context, RecoveryPolicy, CoreDamage) (RecoveryAssessment, error)
}

type RecoveryVerifier interface {
	Preflight(context.Context, RecoveryAssessment) error
	Verify(context.Context, string) error
	Rollback(context.Context, string) error
}
