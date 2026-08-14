package fiber

import "context"

type DamageType string

const (
	DamageCut        DamageType = "cut"
	DamageHighLoss   DamageType = "high-loss"
	DamageJoint      DamageType = "joint"
	DamageSplitter   DamageType = "splitter"
	DamageUnknown    DamageType = "unknown"
)

type RecoveryCandidate struct {
	ID string
	PathID string
	Damage DamageType
	Location Point
	Confidence float64
	AffectedONUs int
	AffectedUsers int
	Priority string
	Reason string
}

type RecoveryPlan struct {
	ID string
	Candidates []RecoveryCandidate
	Steps []string
	ExpectedImpact string
	RollbackPlan string
	RequiresApproval bool
}

// RecoveryPlanner combines topology, OTDR, optical and service evidence into
// an actionable plan. It never performs privileged repair by itself.
type RecoveryPlanner interface {
	Plan(context.Context, FiberTopology, FiberGraph) (RecoveryPlan, error)
	Verify(context.Context, RecoveryPlan) (string, error)
}

type RecoveryExecutor interface {
	ExecuteApproved(context.Context, RecoveryPlan) error
	Rollback(context.Context, RecoveryPlan) error
}
