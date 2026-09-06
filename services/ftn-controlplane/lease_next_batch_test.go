package controlplane

import (
	"testing"
	"time"
)

func TestMemoryLeaseFencingAndExpiry(t *testing.T) {
	s := NewMemoryLeaseStore()
	now := time.Unix(100, 0).UTC()
	first, err := s.Acquire("r1", "owner-a", time.Minute, now)
	if err != nil { t.Fatal(err) }
	if err := s.Validate("r1", "owner-a", first.Fence, now.Add(time.Second)); err != nil { t.Fatal(err) }
	if _, err := s.Renew("r1", "owner-b", first.Fence, time.Minute, now); err == nil { t.Fatal("wrong owner must not renew") }
	if err := s.Validate("r1", "owner-a", first.Fence, now.Add(2*time.Minute)); err == nil { t.Fatal("expired lease must fail validation") }
	second, err := s.Acquire("r1", "owner-b", time.Minute, now.Add(2*time.Minute))
	if err != nil { t.Fatal(err) }
	if second.Fence <= first.Fence { t.Fatal("fence must increase after takeover") }
}
