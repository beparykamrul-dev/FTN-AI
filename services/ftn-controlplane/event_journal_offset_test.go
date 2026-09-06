package controlplane

import "testing"

func TestMemoryEventJournalOffsetCannotMoveBackwards(t *testing.T) {
 j:=NewMemoryEventJournal(); _,_=j.Append(JournalEvent{ID:"e1",TenantID:"t1",Type:"x"})
 if err:=j.CommitOffset("c1","t1",1); err!=nil {t.Fatal(err)}
 if err:=j.CommitOffset("c1","t1",0); err==nil {t.Fatal("expected backwards offset rejection")}
 if got:=j.Offset("c1","t1"); got!=1 {t.Fatalf("offset=%d",got)}
}
