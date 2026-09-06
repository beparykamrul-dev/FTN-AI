package controlplane

import "testing"

func TestMemoryLeaseStoreRejectsWrongOwnerRelease(t *testing.T) {
 s:=NewMemoryLeaseStore(); l,err:=s.Acquire("r","a",1e9,zeroTime()); if err!=nil {t.Fatal(err)}
 if err:=s.Release("r","b",l.Fence); err!=ErrLeaseLost {t.Fatalf("got %v",err)}
}
func zeroTime() (t time.Time) { return }
