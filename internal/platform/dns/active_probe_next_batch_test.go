package dns

import (
 "context"
 "testing"
 "time"
)

func TestDNSProbeValidationAndNilContext(t *testing.T) {
 p := DNSProbe{Address: "127.0.0.1:53", Name: "example.com", Timeout: time.Second}
 if err := p.Validate(); err != nil { t.Fatalf("valid probe rejected: %v", err) }
 if err := (DNSProbe{Address: "127.0.0.1", Name: "example.com", Timeout: time.Second}).Validate(); err == nil { t.Fatal("address without port accepted") }
 r := p.Probe(nil); if r.Error != "context is required" { t.Fatalf("nil context error=%q", r.Error) }
 if got := p.Probe(context.Background()); got.Reachable && got.Error != "" { t.Fatal("unexpected contradictory probe result") }
}
