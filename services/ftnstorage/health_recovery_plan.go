package ftnstorage

// HealthRecoveryPlan is the declarative result of storage-health evaluation.
type HealthRecoveryPlan struct {
	NodeID       string
	State        HealthControllerState
	NeedsRepair  bool
	Quarantine   bool
	Reason       string
}

func BuildHealthRecoveryPlan(nodeID string, h HealthController, p HealthPolicy) HealthRecoveryPlan {
	return HealthRecoveryPlan{
		NodeID:      nodeID,
		State:       h.State,
		NeedsRepair: h.NeedsRecovery(),
		Quarantine:  p.ShouldQuarantine(h.Failures),
		Reason:      "continuous-health-evaluation",
	}
}
