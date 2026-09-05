package dns

import "testing"
func TestDMARCRejectsTooManyReportURIs(t *testing.T){r:=DMARCRecord{Domain:"example.com",Policy:DMARCReject,AggregateReportURI:make([]string,33)};if err:=r.Validate();err==nil{t.Fatal("expected report URI bound")}}
func TestDMARCRejectsInvalidPercentage(t *testing.T){r:=DMARCRecord{Domain:"example.com",Policy:DMARCReject,Percentage:101};if err:=r.Validate();err==nil{t.Fatal("expected percentage rejection")}}
