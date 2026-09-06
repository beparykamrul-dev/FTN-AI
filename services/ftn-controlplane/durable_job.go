package controlplane

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type JobState string

const (
	JobPending   JobState = "pending"
	JobRunning   JobState = "running"
	JobSucceeded JobState = "succeeded"
	JobFailed    JobState = "failed"
	JobCancelled JobState = "cancelled"
)

var (
	ErrJobNotFound      = errors.New("job not found")
	ErrInvalidJobState  = errors.New("invalid job state transition")
	ErrDuplicateJob     = errors.New("duplicate idempotency key")
	ErrImmutableJob     = errors.New("job identity is immutable")
)

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

type JobStore interface {
	Create(DurableJob) error
	Get(string) (DurableJob, error)
	Update(DurableJob) error
}

type MemoryJobStore struct {
	mu   sync.RWMutex
	jobs map[string]DurableJob
	keys map[string]string
}

func NewMemoryJobStore() *MemoryJobStore {
	return &MemoryJobStore{jobs: make(map[string]DurableJob), keys: make(map[string]string)}
}

func (s *MemoryJobStore) Create(job DurableJob) error {
	if s == nil {
		return ErrJobNotFound
	}
	job.ID = strings.TrimSpace(job.ID)
	job.TenantID = strings.TrimSpace(job.TenantID)
	job.IdempotencyKey = strings.TrimSpace(job.IdempotencyKey)
	if job.ID == "" || job.TenantID == "" {
		return errors.New("job id and tenant id are required")
	}
	if job.State == "" {
		job.State = JobPending
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	} else {
		job.CreatedAt = job.CreatedAt.UTC()
	}
	job.UpdatedAt = job.CreatedAt
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.jobs == nil {
		s.jobs = make(map[string]DurableJob)
	}
	if s.keys == nil {
		s.keys = make(map[string]string)
	}
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
	if s == nil {
		return DurableJob{}, ErrJobNotFound
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[strings.TrimSpace(id)]
	if !ok {
		return DurableJob{}, ErrJobNotFound
	}
	return job, nil
}

func (s *MemoryJobStore) Update(job DurableJob) error {
	if s == nil {
		return ErrJobNotFound
	}
	job.ID = strings.TrimSpace(job.ID)
	job.TenantID = strings.TrimSpace(job.TenantID)
	job.IdempotencyKey = strings.TrimSpace(job.IdempotencyKey)
	if job.ID == "" {
		return ErrJobNotFound
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.jobs[job.ID]
	if !ok {
		return ErrJobNotFound
	}
	if job.TenantID != current.TenantID || job.IdempotencyKey != current.IdempotencyKey {
		return ErrImmutableJob
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = current.CreatedAt
	} else {
		job.CreatedAt = current.CreatedAt
	}
	if job.UpdatedAt.IsZero() {
		job.UpdatedAt = time.Now().UTC()
	} else {
		job.UpdatedAt = job.UpdatedAt.UTC()
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
	if now.IsZero() {
		now = time.Now().UTC()
	} else {
		now = now.UTC()
	}
	job.State = next
	job.UpdatedAt = now
	if next == JobRunning {
		job.Attempt++
	}
	return job, nil
}

func CheckpointJob(job DurableJob, checkpoint string, now time.Time) DurableJob {
	if now.IsZero() {
		now = time.Now().UTC()
	} else {
		now = now.UTC()
	}
	job.Checkpoint = strings.TrimSpace(checkpoint)
	job.UpdatedAt = now
	return job
}
