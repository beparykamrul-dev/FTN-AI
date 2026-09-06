package dns

import "testing"

type nextProviderAdapter struct{}
func (nextProviderAdapter) Type() ProviderType { return ProviderPowerDNS }
func (nextProviderAdapter) ApplyZone(_ context.Context, _ Zone) error { return nil }
func (nextProviderAdapter) DeleteZone(_ context.Context, _ string) error { return nil }

func TestProviderEndpointNormalizationRejectsCRLF(t *testing.T) {
	if got := NormalizeProviderEndpoint("https://dns.example/api\r\n"); got != "" { t.Fatalf("unsafe endpoint normalized to %q", got) }
	if got := NormalizeProviderEndpoint("https://dns.example/api/"); got != "https://dns.example/api" { t.Fatalf("endpoint=%q", got) }
}
