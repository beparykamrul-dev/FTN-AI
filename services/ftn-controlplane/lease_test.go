package controlplane

import (
	"testing"
	"time"
)

func TestFencedLeasePreventsStaleOwner(t *testing.T) {
	now := time.Unix(100, 0)
	s := NewMemoryLeaseStore()
	first, err := s.Acquire("job-1", "worker-a", time.Second, now)
	if err != nil || first.Fence != 1 {
		t.Fatalf("first acquire: %+v %v", first, err)
	}
	if _, err := s.Acquire("job-1", "worker-b", time.Second, now.Add(500*time.Millisecond)); err != ErrLeaseLost {
		t.Fatalf("expected competing owner rejection, got %v", err)
	}
	if err := s.Validate("job-1", "worker-a", first.Fence, now.Add(1500*time.Millisecond)); err != ErrLeaseLost {
		t.Fatalf("expected expired lease rejection, got %v", err)
	}
	second, err := s.Acquire("job-1", "worker-b", time.Second, now.Add(1500*time.Millisecond))
	if err != nil || second.Fence != 2 {
		t.Fatalf("second acquire: %+v %v", second, err)
	}
	if err := s.Validate("job-1", "worker-a", first.Fence, now.Add(1600*time.Millisecond)); err != ErrLeaseLost {
		t.Fatalf("stale fence accepted: %v", err)
	}
	if err := s.Validate("job-1", "worker-b", second.Fence, now.Add(1600*time.Millisecond)); err != nil {
		t.Fatalf("current fence rejected: %v", err)
	}
}
