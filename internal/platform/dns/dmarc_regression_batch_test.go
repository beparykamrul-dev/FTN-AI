package dns

import "testing"

func TestDMARCRejectsUnsafeReportURI(t *testing.T) {
	r := DMARCRecord{Domain: "example.com", Policy: DMARCReject, AggregateReportURI: []string{"ftp://example.com/report"}}
	if r.Validate() == nil { t.Fatal("unsupported DMARC report URI scheme must be rejected") }
}

func TestDMARCTXTIncludesPolicy(t *testing.T) {
	r := DMARCRecord{Domain: "example.com", Policy: DMARCReject, Percentage: 100}
	got, err := r.TXTValue(); if err != nil { t.Fatal(err) }
	if got == "" { t.Fatal("DMARC TXT value must not be empty") }
}
