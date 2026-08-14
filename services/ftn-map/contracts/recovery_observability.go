package fiber

import "context"

type RecoveryState string
const (
	RecoveryPlanned RecoveryState = "planned"
	RecoveryRunning RecoveryState = "running"
	RecoveryVerified RecoveryState = "verified"
	RecoveryRolledBack RecoveryState = "rolled-back"
	RecoveryBlocked RecoveryState = "blocked"
)

type RecoveryAttempt struct {
	ID string
	CoreID string
	PathID string
	State RecoveryState
	BeforeLossDb float64
	AfterLossDb float64
	BeforeRxPowerDbm float64
	AfterRxPowerDbm float64
	Confidence float64
	StartedAt string
	FinishedAt string
	Reason string
}

type RecoveryRepository interface {
	RecordAttempt(context.Context, RecoveryAttempt) error
	LatestAttempt(context.Context, string) (RecoveryAttempt, error)
}

// Recovery verification is deliberately independent from the AI decision so
// an unhealthy result can trigger rollback instead of being self-certified.
type RecoveryVerifier interface {
	Verify(context.Context, RecoveryAttempt) (bool, string, error)
}
