package controlplane

import "testing"

func TestMemoryEventJournalSequencesAreTenantScoped(t *testing.T) {
 j:=NewMemoryEventJournal(); a,_:=j.Append(JournalEvent{ID:"a",TenantID:"t1",Type:"x"}); b,_:=j.Append(JournalEvent{ID:"b",TenantID:"t2",Type:"x"})
 if a.Sequence!=1 || b.Sequence!=1 {t.Fatalf("sequences=%d,%d",a.Sequence,b.Sequence)}
}
