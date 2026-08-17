package controlplane

import (
	"errors"
	"sync"
	"time"
)

// JobState is the durable lifecycle state of a control-plane job.
type JobState string

const (
	JobPending   JobState = "pending"
	JobRunning   JobState = "running"
	JobSucceeded JobState = "succeeded"
	JobFailed    JobState = "failed"
	JobCancelled JobState = "cancelled"
)

var (
	ErrJobNotFound    = errors.New("job not found")
	ErrInvalidJobState = errors.New("invalid job state transition")
	ErrDuplicateJob   = errors.New("duplicate idempotency key")
)

// DurableJob is the in-memory state contract. A persistent JobStore can replace
// the implementation without changing lifecycle semantics.
type DurableJob struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	State          JobState  `json:"state"`
	Attempt        int       `json:"attempt"`
	Checkpoint     string    `json:"checkpoint"`
	LastError      string    `json:"last_error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// JobStore defines the minimum durable lifecycle operations required by FTN.OS.
type JobStore interface {
	Create(job DurableJob) error
	Get(id string) (DurableJob, error)
	Update(job DurableJob) error
}

// MemoryJobStore is a deterministic reference store for unit/integration tests.
// Production deployments should provide a transactional persistent adapter.
type MemoryJobStore struct {
	mu   sync.RWMutex
	jobs map[string]DurableJob
	keys map[string]string
}

func NewMemoryJobStore() *MemoryJobStore {
	return &MemoryJobStore{jobs: make(map[string]DurableJob), keys: make(map[string]string)}
}

func (s *MemoryJobStore) Create(job DurableJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[job.ID]; ok {
		return ErrDuplicateJob
	}
	if job.IdempotencyKey != "" {
		if _, ok := s.keys[job.IdempotencyKey]; ok {
			return ErrDuplicateJob
		}
		s.keys[job.IdempotencyKey] = job.ID
	}
	s.jobs[job.ID] = job
	return nil
}

func (s *MemoryJobStore) Get(id string) (DurableJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	if !ok {
		return DurableJob{}, ErrJobNotFound
	}
	return job, nil
}

func (s *MemoryJobStore) Update(job DurableJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[job.ID]; !ok {
		return ErrJobNotFound
	}
	s.jobs[job.ID] = job
	return nil
}

func validJobTransition(from, to JobState) bool {
	switch from {
	case JobPending:
		return to == JobRunning || to == JobCancelled
	case JobRunning:
		return to == JobSucceeded || to == JobFailed || to == JobCancelled || to == JobRunning
	default:
		return false
	}
}

func TransitionJob(job DurableJob, next JobState, now time.Time) (DurableJob, error) {
	if !validJobTransition(job.State, next) {
		return DurableJob{}, ErrInvalidJobState
	}
	job.State = next
	job.UpdatedAt = now
	if next == JobRunning {
		job.Attempt++
	}
	return job, nil
}

func CheckpointJob(job DurableJob, checkpoint string, now time.Time) DurableJob {
	job.Checkpoint = checkpoint
	job.UpdatedAt = now
	return job
}
