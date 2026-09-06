package dns

import "testing"

func TestDMARCValidationRejectsUnsafeURI(t *testing.T) {
	r := DMARCRecord{Domain: "example.com", Policy: DMARCReject, Percentage: 100, AggregateReportURI: []string{"mailto:dmarc@example.com"}}
	if err := r.Validate(); err != nil { t.Fatalf("valid DMARC rejected: %v", err) }
	got, err := r.TXTValue(); if err != nil || got != "v=DMARC1; p=reject; pct=100; rua=mailto:dmarc@example.com" { t.Fatalf("TXT=%q err=%v", got, err) }
	if err := (DMARCRecord{Domain: "example.com", Policy: DMARCReject, AggregateReportURI: []string{"http://example.com/report"}}).Validate(); err == nil { t.Fatal("unsupported report URI accepted") }
}
