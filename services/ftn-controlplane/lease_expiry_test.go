package controlplane

import (
 "testing"
 "time"
)

func TestMemoryLeaseStoreExpiresAtBoundary(t *testing.T) {
 s:=NewMemoryLeaseStore(); now:=time.Date(2026,1,1,0,0,0,0,time.UTC); l,_:=s.Acquire("r1","w1",time.Minute,now)
 if err:=s.Validate("r1","w1",l.Fence,now.Add(time.Minute)); err!=ErrLeaseLost {t.Fatalf("got %v",err)}
}
