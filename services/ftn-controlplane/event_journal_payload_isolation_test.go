package controlplane

import "testing"

func TestMemoryEventJournalAppendCopiesPayload(t *testing.T) {
 j:=NewMemoryEventJournal(); payload:=[]byte("abc")
 e,_:=j.Append(JournalEvent{ID:"e1",TenantID:"t1",Type:"x",Payload:payload}); payload[0]='z'
 got,_:=j.ReadAfter("t1",0,1); if string(got[0].Payload)!="abc" {t.Fatalf("payload mutated: %q",got[0].Payload)}
 e.Payload[0]='y'; got,_=j.ReadAfter("t1",0,1); if string(got[0].Payload)!="abc" {t.Fatalf("stored payload leaked: %q",got[0].Payload)}
}
