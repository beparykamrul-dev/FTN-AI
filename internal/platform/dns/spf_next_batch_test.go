package dns

import "testing"

func TestSPFValidationAndTXTBoundary(t *testing.T) {
	r := SPFRecord{Domain: "Example.COM.", Mechanisms: []string{"a"}, IP4: []string{"192.0.2.0/24"}, All: "fail"}
	if err := r.Validate(); err != nil { t.Fatalf("valid SPF rejected: %v", err) }
	got, err := r.TXTValue(); if err != nil || got != "v=spf1 a ip4:192.0.2.0/24 -all" { t.Fatalf("TXT=%q err=%v", got, err) }
	if err := (SPFRecord{Domain: "example.com", All: "invalid"}).Validate(); err == nil { t.Fatal("invalid all qualifier accepted") }
}
