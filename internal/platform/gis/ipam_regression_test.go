package gis

import "testing"

func TestIPAMRejectsMissingID(t *testing.T) {
	if err := NewIPAM().Upsert(IPAsset{IP:"192.0.2.1"}); err == nil { t.Fatal("IP asset without ID must be rejected") }
}

func TestIPAMListIsDeterministic(t *testing.T) {
	i := NewIPAM()
	_ = i.Upsert(IPAsset{ID:"b",IP:"192.0.2.2"})
	_ = i.Upsert(IPAsset{ID:"a",IP:"192.0.2.1"})
	out := i.List()
	if len(out) != 2 || out[0].ID != "a" || out[1].ID != "b" { t.Fatalf("unexpected order: %#v", out) }
}
