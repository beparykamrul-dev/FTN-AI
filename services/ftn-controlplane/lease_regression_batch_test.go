package controlplane

import (
	"testing"
	"time"
)

func TestLeaseAcquireRenewReleaseFence(t *testing.T) {
	s := NewMemoryLeaseStore()
	now := time.Unix(100, 0).UTC()
	lease, err := s.Acquire("resource", "owner", time.Minute, now)
	if err != nil { t.Fatal(err) }
	if err := s.Validate("resource", "owner", lease.Fence, now); err != nil { t.Fatal(err) }
	if _, err := s.Renew("resource", "other", lease.Fence, time.Minute, now); err == nil { t.Fatal("wrong owner must not renew lease") }
	if err := s.Release("resource", "owner", lease.Fence); err != nil { t.Fatal(err) }
	if err := s.Validate("resource", "owner", lease.Fence, now); err == nil { t.Fatal("released lease must no longer validate") }
}
