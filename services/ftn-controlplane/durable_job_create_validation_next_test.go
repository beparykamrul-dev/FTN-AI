package controlplane

import "testing"

func TestMemoryJobStoreRejectsMissingIdentity(t *testing.T) {
 s:=NewMemoryJobStore(); if err:=s.Create(DurableJob{ID:"",TenantID:"t"}); err==nil {t.Fatal("expected id error")}; if err:=s.Create(DurableJob{ID:"j",TenantID:""}); err==nil {t.Fatal("expected tenant error")}
}
