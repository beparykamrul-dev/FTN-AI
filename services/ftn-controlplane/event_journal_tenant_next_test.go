package controlplane

import "testing"

func TestMemoryEventJournalTenantIsolation(t *testing.T) {
 j:=NewMemoryEventJournal(); _,_ = j.Append(JournalEvent{ID:"a",TenantID:"t1",Type:"x"}); _,_ = j.Append(JournalEvent{ID:"b",TenantID:"t2",Type:"x"})
 got,err:=j.ReadAfter("t1",0,10); if err!=nil {t.Fatal(err)}
 if len(got)!=1 || got[0].TenantID!="t1" {t.Fatalf("unexpected tenant results: %+v",got)}
}
