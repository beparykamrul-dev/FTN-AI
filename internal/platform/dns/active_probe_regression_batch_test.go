package dns

import (
	"context"
	"testing"
	"time"
)

func TestDNSProbeValidationRequiresHostPort(t *testing.T) {
	p := DNSProbe{Address: "127.0.0.1", Name: "example.com", Timeout: time.Second}
	if p.Validate() == nil { t.Fatal("DNS probe address without port must be rejected") }
}

func TestDNSProbeNilContextFailsClosed(t *testing.T) {
	p := DNSProbe{Address: "127.0.0.1:53", Name: "example.com", Timeout: time.Second}
	if r := p.Probe(context.TODO()); r.Reachable { t.Fatal("unreachable test endpoint must not be reported reachable") }
}
