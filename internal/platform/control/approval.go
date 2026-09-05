package control

import (
	"errors"
	"sync"
	"time"
)

type ApprovalStatus string

const (
	Pending  ApprovalStatus = "pending"
	Approved ApprovalStatus = "approved"
	Rejected ApprovalStatus = "rejected"
)

type Request struct {
	ID        string         `json:"id"`
	ServerID  string         `json:"server_id"`
	Operation Operation      `json:"operation"`
	Reason    string         `json:"reason"`
	Status    ApprovalStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ApprovalStore struct {
	mu sync.RWMutex
	m  map[string]Request
}

func NewApprovalStore() *ApprovalStore {
	return &ApprovalStore{m: make(map[string]Request)}
}

func (s *ApprovalStore) Create(r Request) error {
	if r.ID == "" || r.ServerID == "" || r.Operation == "" {
		return errors.New("invalid approval request")
	}
	if err := Validate(r.Operation); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.m[r.ID]; exists {
		return errors.New("approval request already exists")
	}
	now := time.Now().UTC()
	r.Status = Pending
	r.CreatedAt = now
	r.UpdatedAt = now
	s.m[r.ID] = r
	return nil
}

func (s *ApprovalStore) SetStatus(id string, status ApprovalStatus) error {
	if status != Approved && status != Rejected {
		return errors.New("invalid approval status")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.m[id]
	if !ok {
		return errors.New("approval request not found")
	}
	r.Status = status
	r.UpdatedAt = time.Now().UTC()
	s.m[id] = r
	return nil
}

func (s *ApprovalStore) Get(id string) (Request, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.m[id]
	return r, ok
}
