package monitoring

import (
	"context"
	"time"
)

type PollJob struct {
	Target SNMPTarget `json:"target"`
	Profile SNMPProfile `json:"profile"`
}

type PollResult struct {
	TargetID string `json:"target_id"`
	Samples []SNMPSample `json:"samples,omitempty"`
	Error string `json:"error,omitempty"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

// PollScheduler provides cancellation-aware scheduling around the injected
// SNMP transport. It never performs configuration writes.
type PollScheduler struct { Poller SNMPPoller }

func (s PollScheduler) RunOnce(ctx context.Context, job PollJob) PollResult {
	started := time.Now().UTC()
	result := PollResult{TargetID: job.Target.ID, StartedAt: started}
	if err := ctx.Err(); err != nil {
		result.Error = err.Error()
		result.FinishedAt = time.Now().UTC()
		return result
	}
	samples, err := PollTarget(s.Poller, job.Target, job.Profile.OIDs)
	if err != nil { result.Error = err.Error() } else { result.Samples = samples }
	result.FinishedAt = time.Now().UTC()
	return result
}
