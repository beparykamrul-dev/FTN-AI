package controlplane

import "testing"

func TestMemoryEventJournalRejectsBackwardOffset(t *testing.T) {
 j:=NewMemoryEventJournal(); _,err:=j.Append(JournalEvent{ID:"e",TenantID:"t",Type:"x"}); if err!=nil {t.Fatal(err)}
 if err:=j.CommitOffset("c","t",1); err!=nil {t.Fatal(err)}
 if err:=j.CommitOffset("c","t",0); err==nil {t.Fatal("expected backward offset rejection")}
}
