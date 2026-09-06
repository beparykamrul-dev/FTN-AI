package dns

import "testing"

func TestSPFRejectsInvalidNetwork(t *testing.T) {
	r := SPFRecord{Domain: "example.com", IP4: []string{"not-a-cidr"}}
	if r.Validate() == nil { t.Fatal("invalid SPF network must be rejected") }
}

func TestSPFTXTIsDeterministic(t *testing.T) {
	r := SPFRecord{Domain: "example.com", Include: []string{"b.example", "a.example"}, All: "fail"}
	a, err := r.TXTValue(); if err != nil { t.Fatal(err) }
	b, err := r.TXTValue(); if err != nil { t.Fatal(err) }
	if a != b { t.Fatal("SPF rendering must be deterministic") }
}
