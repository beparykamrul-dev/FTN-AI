package controlplane

import "testing"

func TestMemoryEventJournalDuplicateAppendKeepsSequence(t *testing.T) {
 j:=NewMemoryEventJournal(); a,err:=j.Append(JournalEvent{ID:"e",TenantID:"t",Type:"x",Payload:[]byte("a")}); if err!=nil {t.Fatal(err)}
 b,err:=j.Append(JournalEvent{ID:"e",TenantID:"t",Type:"x",Payload:[]byte("b")}); if err!=nil {t.Fatal(err)}
 if b.Sequence!=a.Sequence || string(b.Payload)!="a" {t.Fatalf("unexpected duplicate: %+v",b)}
}
