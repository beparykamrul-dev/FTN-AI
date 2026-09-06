package controlplane

import "testing"

func TestMemoryJobStoreRejectsDuplicateID(t *testing.T) {
 s:=NewMemoryJobStore(); j:=DurableJob{ID:"j",TenantID:"t"}; if err:=s.Create(j); err!=nil {t.Fatal(err)}; if err:=s.Create(j); err!=ErrDuplicateJob {t.Fatalf("got %v",err)}
}
