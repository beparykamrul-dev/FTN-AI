package fiber

import "context"

type RecoveryGuardState string

const (
	GuardReady      RecoveryGuardState = "ready"
	GuardBlocked    RecoveryGuardState = "blocked"
	GuardVerifying  RecoveryGuardState = "verifying"
	GuardRollback   RecoveryGuardState = "rollback"
)

type RecoveryGuard struct {
	CoreID string
	State RecoveryGuardState
	BaselineLossDb float64
	CurrentLossDb float64
	BaselineRxPowerDbm float64
	CurrentRxPowerDbm float64
	Confidence float64
	Attempts int
	MaxAttempts int
	Reason string
}

// RecoveryGuard prevents repeated or unsafe automatic changes. Every
// mitigation must establish a baseline, pass preflight, and verify optical
// health before being considered successful.
type RecoveryGuardRepository interface {
	Load(context.Context, string) (RecoveryGuard, error)
	Save(context.Context, RecoveryGuard) error
}

type RecoveryGuardEngine interface {
	Preflight(context.Context, RecoveryAction) (RecoveryGuard, error)
	Verify(context.Context, RecoveryAction) (RecoveryGuard, error)
	Rollback(context.Context, RecoveryAction) error
}
