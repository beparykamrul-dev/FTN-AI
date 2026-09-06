package controlplane

import "testing"

func TestMemoryEventJournalCopiesPayload(t *testing.T) {
 j:=NewMemoryEventJournal(); p:=[]byte("abc"); e,err:=j.Append(JournalEvent{ID:"e",TenantID:"t",Type:"x",Payload:p}); if err!=nil {t.Fatal(err)}
 p[0]='z'; got,err:=j.ReadAfter("t",0,1); if err!=nil {t.Fatal(err)}
 if string(e.Payload)!="abc" || string(got[0].Payload)!="abc" {t.Fatalf("payload alias detected")}
}
