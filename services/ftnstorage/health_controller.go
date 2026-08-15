package ftnstorage

import "time"

// HealthControllerState is the state machine used to keep storage under
// continuous observation and bounded recovery.
type HealthControllerState string

const (
	HealthHealthy   HealthControllerState = "healthy"
	HealthDegraded  HealthControllerState = "degraded"
	HealthRepairing HealthControllerState = "repairing"
	HealthQuarantine HealthControllerState = "quarantine"
)

type HealthController struct {
	State      HealthControllerState
	LastCheck  time.Time
	Failures   uint32
	Repairing  bool
}

func (h *HealthController) Observe(ok bool, now time.Time) {
	h.LastCheck = now.UTC()
	if ok {
		h.Failures = 0
		h.Repairing = false
		h.State = HealthHealthy
		return
	}
	h.Failures++
	h.State = HealthDegraded
}

func (h HealthController) NeedsRecovery() bool {
	return h.State == HealthDegraded || h.State == HealthRepairing
}
