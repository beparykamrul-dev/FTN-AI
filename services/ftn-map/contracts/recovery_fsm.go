package fiber

import "context"

type RecoveryState string
const (
	RecoveryObserved RecoveryState = "observed"
	RecoveryDiagnosing RecoveryState = "diagnosing"
	RecoveryPreflight RecoveryState = "preflight"
	RecoveryMitigating RecoveryState = "mitigating"
	RecoveryVerifying RecoveryState = "verifying"
	RecoveryRecovered RecoveryState = "recovered"
	RecoveryRollback RecoveryState = "rollback"
	RecoveryEscalated RecoveryState = "escalated"
)

type RecoveryIncident struct {
	ID string
	DeviceID string
	CoreID string
	PathID string
	State RecoveryState
	Confidence float64
	Reason string
}

// RecoveryStateMachine provides an explicit, auditable lifecycle for autonomous
// mitigation. Every state transition is observable and policy-controlled.
type RecoveryStateMachine interface {
	Start(context.Context, RecoveryIncident) error
	Transition(context.Context, string, RecoveryState, string) error
	Current(context.Context, string) (RecoveryIncident, error)
	Close(context.Context, string) error
}
