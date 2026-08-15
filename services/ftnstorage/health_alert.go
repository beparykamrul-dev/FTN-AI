package ftnstorage

// HealthAlert is emitted by the storage health controller for operator/telemetry layers.
type HealthAlert struct {
	NodeID   string
	State    HealthControllerState
	Failures uint32
	Reason   string
}

func BuildHealthAlert(nodeID string, h HealthController) HealthAlert {
	return HealthAlert{NodeID: nodeID, State: h.State, Failures: h.Failures, Reason: "storage-health-state-change"}
}
