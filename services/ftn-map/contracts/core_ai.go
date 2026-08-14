package fiber

import "context"

// FiberCoreAI models one addressable optical core as an independently observed
// asset. It does not imply that software can physically repair glass; it can
// monitor, correlate, protect service and recommend authorized recovery.
type FiberCoreState struct {
	PathID string
	Core int
	State string
	TxPowerDbm float64
	RxPowerDbm float64
	LossDb float64
	DistanceMeters float64
	OTDRStatus string
	HealthScore float64
	LastSeen string
}

type CoreFinding struct {
	PathID string
	Core int
	Finding string
	Severity string
	Confidence float64
	AffectedONUs int
	AffectedUsers int
}

type CoreRecoveryPlan struct {
	PathID string
	Core int
	Actions []string
	RequiresFieldRepair bool
	Verified bool
}

// CoreAI provides per-core monitoring, diagnosis and recovery planning.
type CoreAI interface {
	Observe(context.Context, FiberCoreState) error
	Analyze(context.Context, FiberCoreState) (CoreFinding, error)
	PlanRecovery(context.Context, FiberCoreState, CoreFinding) (CoreRecoveryPlan, error)
	Verify(context.Context, FiberCoreState) (bool, error)
}
