package fiber

import "context"

type DamageType string
const (
	DamageHighLoss DamageType = "high-loss"
	DamageDirtyConnector DamageType = "dirty-connector"
	DamageBadSplice DamageType = "bad-splice"
	DamageBend DamageType = "excessive-bend"
	DamageWeakOpticalSignal DamageType = "weak-optical-signal"
	DamageSplitter DamageType = "splitter-degradation"
	DamagePort DamageType = "port-degradation"
)

type RecoveryAction string
const (
	ActionRebalanceOpticalPower RecoveryAction = "rebalance-optical-power"
	ActionRerouteSpareCore RecoveryAction = "reroute-spare-core"
	ActionSwitchSparePort RecoveryAction = "switch-spare-port"
	ActionFlagFieldRepair RecoveryAction = "flag-field-repair"
	ActionReconfigureSplitter RecoveryAction = "reconfigure-splitter"
	ActionCleanConnector RecoveryAction = "clean-connector"
)

type NonCutDamage struct {
	ID string
	PathID string
	EntityID string
	Type DamageType
	Severity string
	LossDb float64
	Confidence float64
	DetectedAt string
}

type RecoveryCandidate struct {
	DamageID string
	Action RecoveryAction
	ExpectedLossDb float64
	AffectedUsers int
	RequiresFieldWork bool
	RequiresApproval bool
	Reason string
}

// NonCutRecovery handles degradations that may be mitigated without physically
// cutting/re-splicing the fiber. It never claims a damaged physical fiber is
// magically repaired; field repair remains the fallback when attenuation or
// physical damage cannot be safely mitigated remotely.
type NonCutRecovery interface {
	Detect(context.Context, FiberTopology) ([]NonCutDamage, error)
	Plan(context.Context, NonCutDamage) ([]RecoveryCandidate, error)
	Verify(context.Context, RecoveryCandidate) (bool, error)
}
