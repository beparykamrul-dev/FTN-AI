package controlplane

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLeaseNotFound = errors.New("lease not found")
	ErrLeaseLost    = errors.New("lease lost")
)

type Lease struct {
	ResourceID string    `json:"resource_id"`
	OwnerID    string    `json:"owner_id"`
	Fence      uint64    `json:"fence"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// LeaseStore provides fencing-token based ownership for single-writer jobs.
type LeaseStore interface {
	Acquire(resourceID, ownerID string, ttl time.Duration, now time.Time) (Lease, error)
	Renew(resourceID, ownerID string, fence uint64, ttl time.Duration, now time.Time) (Lease, error)
	Release(resourceID, ownerID string, fence uint64) error
	Validate(resourceID, ownerID string, fence uint64, now time.Time) error
}

type MemoryLeaseStore struct {
	mu     sync.Mutex
	leases map[string]Lease
	next   map[string]uint64
}

func NewMemoryLeaseStore() *MemoryLeaseStore {
	return &MemoryLeaseStore{leases: make(map[string]Lease), next: make(map[string]uint64)}
}

func (s *MemoryLeaseStore) Acquire(resourceID, ownerID string, ttl time.Duration, now time.Time) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, ok := s.leases[resourceID]; ok && current.ExpiresAt.After(now) {
		if current.OwnerID == ownerID {
			return current, nil
		}
		return Lease{}, ErrLeaseLost
	}
	s.next[resourceID]++
	lease := Lease{ResourceID: resourceID, OwnerID: ownerID, Fence: s.next[resourceID], ExpiresAt: now.Add(ttl)}
	s.leases[resourceID] = lease
	return lease, nil
}

func (s *MemoryLeaseStore) Renew(resourceID, ownerID string, fence uint64, ttl time.Duration, now time.Time) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.leases[resourceID]
	if !ok || current.OwnerID != ownerID || current.Fence != fence || !current.ExpiresAt.After(now) {
		return Lease{}, ErrLeaseLost
	}
	current.ExpiresAt = now.Add(ttl)
	s.leases[resourceID] = current
	return current, nil
}

func (s *MemoryLeaseStore) Release(resourceID, ownerID string, fence uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.leases[resourceID]
	if !ok {
		return ErrLeaseNotFound
	}
	if current.OwnerID != ownerID || current.Fence != fence {
		return ErrLeaseLost
	}
	delete(s.leases, resourceID)
	return nil
}

func (s *MemoryLeaseStore) Validate(resourceID, ownerID string, fence uint64, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.leases[resourceID]
	if !ok || current.OwnerID != ownerID || current.Fence != fence || !current.ExpiresAt.After(now) {
		return ErrLeaseLost
	}
	return nil
}
