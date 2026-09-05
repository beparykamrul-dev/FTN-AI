package dns

import "testing"

func TestSPFRejectsInvalidQualifier(t *testing.T) {
	if err := (SPFRecord{Domain:"example.com", All:"bogus"}).Validate(); err == nil { t.Fatal("invalid SPF qualifier must be rejected") }
}

func TestSPFSortsIncludesDeterministically(t *testing.T) {
	r := SPFRecord{Domain:"example.com", Include:[]string{"b.example","a.example"}, All:"fail"}
	got, err := r.TXTValue()
	if err != nil || got != "v=spf1 include:a.example include:b.example -all" { t.Fatalf("unexpected SPF TXT: %q err=%v", got, err) }
}
