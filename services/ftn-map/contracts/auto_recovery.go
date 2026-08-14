package fiber

import "context"

type DamageBand string
const (
	DamageNormal DamageBand = "normal"
	DamageMinor DamageBand = "minor"
	DamageModerate DamageBand = "moderate"
	DamageCritical DamageBand = "critical"
	DamageCut DamageBand = "cut"
	DamageBrokenGlass DamageBand = "broken-glass"
)

type CoreDamage struct {
	CoreID string
	PathID string
	DamagePercent float64
	Band DamageBand
	OpticalLossDb float64
	ServiceImpact float64
	Confidence float64
}

type RecoveryAction struct {
	CoreID string
	Action string
	ExpectedImprovementPercent float64
	RequiresFieldRepair bool
	RequiresApproval bool
	Reason string
}

type AutoRecoveryEngine interface {
	Assess(context.Context, CoreDamage) (RecoveryAction, error)
	Apply(context.Context, RecoveryAction) error
	Verify(context.Context, string) error
}
