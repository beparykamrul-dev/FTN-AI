package controlplane

import "testing"

func TestMemoryEventJournalRejectsMissingIdentity(t *testing.T) {
 j:=NewMemoryEventJournal()
 if _,err:=j.Append(JournalEvent{TenantID:"t",Type:"x"}); err==nil {t.Fatal("expected missing id error")}
 if _,err:=j.Append(JournalEvent{ID:"e",Type:"x"}); err==nil {t.Fatal("expected missing tenant error")}
 if _,err:=j.Append(JournalEvent{ID:"e",TenantID:"t"}); err==nil {t.Fatal("expected missing type error")}
}
