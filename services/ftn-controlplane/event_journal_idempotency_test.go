package controlplane

import "testing"

func TestMemoryEventJournalAppendIsIdempotent(t *testing.T) {
 j:=NewMemoryEventJournal(); e:=JournalEvent{ID:"e1",TenantID:"t1",Type:"job.done",Payload:[]byte("x")}
 first,err:=j.Append(e); if err!=nil {t.Fatal(err)}
 second,err:=j.Append(e); if err!=nil {t.Fatal(err)}
 if first.Sequence!=second.Sequence {t.Fatalf("sequence changed: %d != %d",first.Sequence,second.Sequence)}
 got,_:=j.ReadAfter("t1",0,10); if len(got)!=1 {t.Fatalf("events=%d",len(got))}
}
