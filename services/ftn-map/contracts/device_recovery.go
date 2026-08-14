package fiber

import "context"

type DeviceHealthState string
const (
	DeviceOnline DeviceHealthState = "online"
	DeviceOffline DeviceHealthState = "offline"
	DevicePowerUnknown DeviceHealthState = "power-unknown"
	DeviceRecoveryCandidate DeviceHealthState = "recovery-candidate"
	DeviceRecoveryBlocked DeviceHealthState = "recovery-blocked"
)

type DeviceOutage struct {
	DeviceID string
	StartedAt string
	DurationSeconds int64
	PowerIssue bool
	Health DeviceHealthState
	LastSeen string
	Reason string
}

type DeviceRecoveryPlan struct {
	DeviceID string
	Reason string
	Actions []string
	RequiresApproval bool
	RequiresFieldRepair bool
	ExpectedImpact string
}

// DeviceRecoveryEngine correlates telemetry, reachability, topology and power
// evidence before recommending or applying a recovery action.
type DeviceRecoveryEngine interface {
	Evaluate(context.Context, DiscoveredDevice, DeviceOutage) (DeviceRecoveryPlan, error)
	Apply(context.Context, DeviceRecoveryPlan) error
	Verify(context.Context, string) error
}
