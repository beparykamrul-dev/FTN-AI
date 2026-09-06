package controlplane

import (
 "testing"
 "time"
)

func TestMemoryLeaseStoreRejectsInvalidAcquire(t *testing.T) {
 s := NewMemoryLeaseStore()
 if _, err := s.Acquire("", "owner", time.Minute, time.Now()); err != ErrLeaseLost { t.Fatalf("got %v", err) }
 if _, err := s.Acquire("resource", "", time.Minute, time.Now()); err != ErrLeaseLost { t.Fatalf("got %v", err) }
 if _, err := s.Acquire("resource", "owner", 0, time.Now()); err != ErrLeaseLost { t.Fatalf("got %v", err) }
}
