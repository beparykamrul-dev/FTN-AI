package fiber

import "context"

type RecoveryDecision struct {
	DeviceID string
	Decision string
	Confidence float64
	Evidence []string
	Actions []string
	RequiresApproval bool
	BlockedReason string
}

// DecisionEngine correlates independent evidence before recovery. A single
// stale/missing signal must not be treated as proof of a fault.
type DecisionEngine interface {
	Decide(context.Context, DiscoveredDevice, DeviceOutage, FiberTopology) (RecoveryDecision, error)
}
