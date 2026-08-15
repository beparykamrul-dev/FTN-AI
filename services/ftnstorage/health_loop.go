package ftnstorage

import "time"

// HealthLoopResult is the bounded decision produced after a storage health cycle.
type HealthLoopResult struct {
	State          HealthControllerState
	RepairRequired bool
	Rollback       bool
	Quarantine     bool
	NextCheck      time.Time
}

func RunHealthCycle(h *HealthController, p HealthPolicy, ok bool, now time.Time) HealthLoopResult {
	h.Observe(ok, now)
	return HealthLoopResult{
		State:          h.State,
		RepairRequired: h.NeedsRecovery(),
		Quarantine:     p.ShouldQuarantine(h.Failures),
		NextCheck:      now.UTC().Add(time.Minute),
	}
}
