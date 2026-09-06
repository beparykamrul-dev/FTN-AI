package controlplane

import (
 "testing"
 "time"
)

func TestMemoryLeaseStoreRejectsWrongOwnerRelease(t *testing.T) {
 s:=NewMemoryLeaseStore(); l,err:=s.Acquire("r","a",time.Second,time.Time{}); if err!=nil {t.Fatal(err)}
 if err:=s.Release("r","b",l.Fence); err!=ErrLeaseLost {t.Fatalf("got %v",err)}
}
