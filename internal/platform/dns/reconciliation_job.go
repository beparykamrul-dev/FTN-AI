package dns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ReconciliationJobStatus string

const (
	JobPending ReconciliationJobStatus = "pending"
	JobRunning ReconciliationJobStatus = "running"
	JobSucceeded ReconciliationJobStatus = "succeeded"
	JobFailed ReconciliationJobStatus = "failed"
)

type ReconciliationJob struct {
	Key string `json:"key"`
	Status ReconciliationJobStatus `json:"status"`
	Attempts uint32 `json:"attempts"`
	NextRunAt time.Time `json:"next_run_at"`
	LastError string `json:"last_error,omitempty"`
}

// ReconciliationJobStore is a persistence boundary. Production deployments
// can implement it with PostgreSQL/Redis while the orchestration layer remains
// independent of the storage engine.
type ReconciliationJobStore interface {
	Get(ctx context.Context, key string) (ReconciliationJob, bool, error)
	Put(ctx context.Context, job ReconciliationJob) error
}

type MemoryReconciliationJobStore struct {
	mu sync.RWMutex
	jobs map[string]ReconciliationJob
}

func NewMemoryReconciliationJobStore() *MemoryReconciliationJobStore {
	return &MemoryReconciliationJobStore{jobs: make(map[string]ReconciliationJob)}
}

func (s *MemoryReconciliationJobStore) Get(ctx context.Context, key string) (ReconciliationJob, bool, error) {
	select { case <-ctx.Done(): return ReconciliationJob{}, false, ctx.Err(); default: }
	s.mu.RLock(); defer s.mu.RUnlock()
	job, ok := s.jobs[key]
	return job, ok, nil
}

func (s *MemoryReconciliationJobStore) Put(ctx context.Context, job ReconciliationJob) error {
	if job.Key == "" { return fmt.Errorf("job key is required") }
	select { case <-ctx.Done(): return ctx.Err(); default: }
	s.mu.Lock(); defer s.mu.Unlock()
	s.jobs[job.Key] = job
	return nil
}
