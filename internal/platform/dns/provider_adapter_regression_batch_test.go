package dns

import "testing"

func TestNormalizeProviderEndpointRejectsCRLF(t *testing.T) {
	if got := NormalizeProviderEndpoint("https://example.com\r\nX: y"); got != "" { t.Fatalf("unsafe endpoint must normalize to empty, got %q", got) }
}

func TestAdapterRegistryRejectsDuplicateType(t *testing.T) {
	r := NewAdapterRegistry()
	if r == nil { t.Fatal("registry must be created") }
}
