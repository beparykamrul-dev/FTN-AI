package controlplane

import "testing"

func TestMemoryJobStoreProtectsJobIdentity(t *testing.T) {
 s:=NewMemoryJobStore(); if err:=s.Create(DurableJob{ID:"j",TenantID:"t",IdempotencyKey:"k"}); err!=nil {t.Fatal(err)}
 if err:=s.Update(DurableJob{ID:"j",TenantID:"other",IdempotencyKey:"k",State:JobRunning}); err!=ErrImmutableJob {t.Fatalf("got %v",err)}
 if err:=s.Update(DurableJob{ID:"j",TenantID:"t",IdempotencyKey:"other",State:JobRunning}); err!=ErrImmutableJob {t.Fatalf("got %v",err)}
}
