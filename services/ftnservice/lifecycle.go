package ftnservice

// Lifecycle is the controlled service lifecycle state.
type Lifecycle string

const (
	StateStarting Lifecycle = "starting"
	StateReady    Lifecycle = "ready"
	StateDegraded Lifecycle = "degraded"
	StateStopping Lifecycle = "stopping"
	StateStopped  Lifecycle = "stopped"
)

func CanServe(s Lifecycle) bool { return s == StateReady }
