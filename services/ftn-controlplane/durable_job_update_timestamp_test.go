package controlplane

import (
 "testing"
 "time"
)

func TestMemoryJobStoreUpdatePreservesCreatedAt(t *testing.T) {
 s:=NewMemoryJobStore(); created:=time.Date(2026,1,1,0,0,0,0,time.UTC)
 if err:=s.Create(DurableJob{ID:"j1",TenantID:"t1",CreatedAt:created}); err!=nil {t.Fatal(err)}
 job,_:=s.Get("j1"); job.Checkpoint="x"; job.CreatedAt=time.Time{}
 if err:=s.Update(job); err!=nil {t.Fatal(err)}
 got,_:=s.Get("j1"); if !got.CreatedAt.Equal(created) {t.Fatalf("created_at changed: %v",got.CreatedAt)}
}
