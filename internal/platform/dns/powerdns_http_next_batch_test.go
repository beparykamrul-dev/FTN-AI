package dns

import "testing"

func TestNewPowerDNSHTTPClientNormalizesEndpoint(t *testing.T) {
	c := NewPowerDNSHTTPClient(" https://dns.example/api/ ", "server-1", "secret")
	if c == nil || c.Client == nil { t.Fatal("client must be initialized") }
	if c.Endpoint != "https://dns.example/api" { t.Fatalf("endpoint=%q", c.Endpoint) }
}
