package controlplane

import "testing"

func TestMemoryJobStoreUpdateRejectsIdentityChange(t *testing.T) {
	s := NewMemoryJobStore()
	if err := s.Create(DurableJob{ID: "j1", TenantID: "t1", IdempotencyKey: "k1"}); err != nil { t.Fatal(err) }
	job, _ := s.Get("j1")
	job.TenantID = "t2"
	if err := s.Update(job); err != ErrImmutableJob { t.Fatalf("got %v", err) }
}
