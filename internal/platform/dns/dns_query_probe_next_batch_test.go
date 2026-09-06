package dns

import (
 "context"
 "testing"
 "time"
)

func TestDNSQueryProbeBoundaries(t *testing.T) {
 p := DNSQueryProbe{Address: "127.0.0.1:53", Name: "example.com.", Timeout: time.Second}
 if err := p.Validate(); err != nil { t.Fatalf("valid query probe rejected: %v", err) }
 if err := (DNSQueryProbe{Address: "127.0.0.1:53", Name: "", Timeout: time.Second}).Validate(); err == nil { t.Fatal("empty query name accepted") }
 r := p.Probe(nil); if r.Error != "context is required" { t.Fatalf("nil context error=%q", r.Error) }
 ctx, cancel := context.WithCancel(context.Background()); cancel(); r = p.Probe(ctx); if r.Error == "" { t.Fatal("cancelled context did not fail") }
}
