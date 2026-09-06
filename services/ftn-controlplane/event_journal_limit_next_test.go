package controlplane

import "testing"

func TestMemoryEventJournalReadAfterZeroLimit(t *testing.T) {
 j:=NewMemoryEventJournal(); if _,err:=j.ReadAfter("tenant",0,0); err!=nil {t.Fatal(err)}
}
