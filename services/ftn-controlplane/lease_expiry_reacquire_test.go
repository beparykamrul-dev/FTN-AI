package controlplane

import (
 "testing"
 "time"
)

func TestMemoryLeaseStoreReacquiresAfterExpiry(t *testing.T) {
 s := NewMemoryLeaseStore(); now := time.Unix(100,0).UTC()
 a, err := s.Acquire("r", "a", time.Second, now); if err != nil { t.Fatal(err) }
 b, err := s.Acquire("r", "b", time.Second, now.Add(2*time.Second)); if err != nil { t.Fatal(err) }
 if b.OwnerID != "b" || b.Fence <= a.Fence { t.Fatalf("unexpected lease: %+v", b) }
}
