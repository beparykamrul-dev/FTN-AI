package observability

// BackendRecoveryState tracks whether a degraded telemetry backend can re-enter service.
type BackendRecoveryState struct {
	Name string
	Healthy bool
	ConsecutiveHealthy uint32
	RequiredHealthy uint32
}

func (s BackendRecoveryState) Recovered() bool {
	return s.Name != "" && s.Healthy && s.RequiredHealthy > 0 && s.ConsecutiveHealthy >= s.RequiredHealthy
}
