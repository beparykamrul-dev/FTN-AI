package controlplane

import (
 "testing"
 "time"
)

func TestMemoryLeaseStoreRejectsStaleFence(t *testing.T) {
 s:=NewMemoryLeaseStore(); now:=time.Date(2026,1,1,0,0,0,0,time.UTC)
 a,err:=s.Acquire("r1","w1",time.Minute,now); if err!=nil {t.Fatal(err)}
 if err:=s.Validate("r1","w1",a.Fence+1,now); err!=ErrLeaseLost {t.Fatalf("got %v",err)}
}
