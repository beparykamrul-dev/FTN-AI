package controlplane

import "testing"

func TestMemoryJobStoreRejectsDuplicateIdempotencyKey(t *testing.T) {
 s:=NewMemoryJobStore(); if err:=s.Create(DurableJob{ID:"a",TenantID:"t",IdempotencyKey:"k"}); err!=nil {t.Fatal(err)}; if err:=s.Create(DurableJob{ID:"b",TenantID:"t",IdempotencyKey:"k"}); err!=ErrDuplicateJob {t.Fatalf("got %v",err)}
}
