package model

import "testing"

func TestIncidentNormalizeDeduplicatesEvidence(t *testing.T) {
	i := (Incident{ID:" inc-1 ", Severity:" high ", EvidenceIDs:[]string{" a ","a","","b"}}).Normalize()
	if i.ID != "inc-1" || i.Severity != "high" { t.Fatalf("unexpected normalized identity: %#v", i) }
	if len(i.EvidenceIDs) != 2 || i.EvidenceIDs[0] != "a" || i.EvidenceIDs[1] != "b" { t.Fatalf("unexpected evidence ids: %#v", i.EvidenceIDs) }
}

func TestIncidentValidRejectsMissingIdentity(t *testing.T) {
	if (Incident{Severity:"high"}).Valid() { t.Fatal("incident without id must be invalid") }
}
