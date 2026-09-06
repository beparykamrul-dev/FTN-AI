package mesh

import "testing"

func TestDefaultProtocolCapabilitiesAreComplete(t *testing.T) {
	caps := DefaultProtocolCapabilities()
	if len(caps) < 6 {
		t.Fatalf("expected default protocol capabilities, got %d", len(caps))
	}
	for _, c := range caps {
		if c.Protocol == "" || len(c.Capabilities) == 0 {
			t.Fatalf("invalid protocol capability: %+v", c)
		}
	}
}
