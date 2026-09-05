package mesh

import "testing"

func TestDefaultProtocolCapabilitiesIsNonEmptyAndUnique(t *testing.T) {
	caps := DefaultProtocolCapabilities()
	if len(caps) == 0 { t.Fatal("protocol capability registry must not be empty") }
	seen := map[Protocol]bool{}
	for _, c := range caps {
		if c.Protocol == "" { t.Fatal("protocol capability must have an identity") }
		if seen[c.Protocol] { t.Fatalf("duplicate protocol %q", c.Protocol) }
		seen[c.Protocol] = true
	}
}
