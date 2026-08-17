package controlplane

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrLeaseRequired = errors.New("valid lease required")

type WorkerResult struct {
	Job      DurableJob
	Event    JournalEvent
	Executed bool
}

// IdempotentWorker coordinates a job using the journal and fenced lease.
type IdempotentWorker struct {
	Jobs   JobStore
	Events EventJournal
	Leases LeaseStore
	mu     sync.Mutex
}

func (w *IdempotentWorker) Execute(jobID, workerID, leaseKey string, token uint64, eventType string, now time.Time) (WorkerResult, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	job, err := w.Jobs.Get(jobID)
	if err != nil {
		return WorkerResult{}, err
	}
	if err = w.Leases.Validate(leaseKey, workerID, token, now); err != nil {
		return WorkerResult{}, ErrLeaseRequired
	}
	if job.State == JobSucceeded {
		return WorkerResult{Job: job, Executed: false}, nil
	}
	if job.State == JobPending {
		job, err = TransitionJob(job, JobRunning, now)
		if err != nil {
			return WorkerResult{}, err
		}
		if err = w.Jobs.Update(job); err != nil {
			return WorkerResult{}, err
		}
	}

	event, err := w.Events.Append(JournalEvent{
		ID: fmt.Sprintf("job:%s:attempt:%d", job.ID, job.Attempt),
		TenantID: job.TenantID,
		Type: eventType,
		CorrelationID: job.ID,
		AggregateID: job.ID,
	})
	if err != nil {
		return WorkerResult{}, err
	}
	job, err = TransitionJob(job, JobSucceeded, now)
	if err != nil {
		return WorkerResult{}, err
	}
	if err = w.Jobs.Update(job); err != nil {
		return WorkerResult{}, err
	}
	return WorkerResult{Job: job, Event: event, Executed: true}, nil
}

// ReconcilePending replays pending jobs through the same fenced execution path.
func (w *IdempotentWorker) ReconcilePending(tenantID, workerID, leaseKey string, token uint64, now time.Time, ids []string) error {
	for _, id := range ids {
		job, err := w.Jobs.Get(id)
		if err != nil || job.TenantID != tenantID || job.State != JobPending {
			continue
		}
		if _, err = w.Execute(id, workerID, leaseKey, token, "job.reconciled", now); err != nil {
			return err
		}
	}
	return nil
}
