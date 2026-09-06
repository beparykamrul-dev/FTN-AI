package dns

import (
	"context"
	"testing"
	"time"
)

func TestDNSProbeValidationBounds(t *testing.T) {
	if err := (DNSProbe{Address:"127.0.0.1", Name:"x", Timeout:time.Second}).Validate(); err == nil { t.Fatal("host without port must fail") }
	if err := (DNSProbe{Address:"127.0.0.1:53", Timeout:31*time.Second}).Validate(); err == nil { t.Fatal("timeout over 30s must fail") }
	r := (DNSProbe{Address:"127.0.0.1:53", Timeout:time.Second}).Probe(context.Background())
	if r.Address == "" { t.Fatal("probe must preserve address") }
}
