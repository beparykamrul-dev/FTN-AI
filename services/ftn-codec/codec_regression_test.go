package codec

import "testing"

func TestCapabilityValidationRejectsMissingIdentity(t *testing.T) {
	if (Capability{Class:"network"}).Valid() { t.Fatal("capability without id must be invalid") }
}

func TestDefaultCapabilitiesAreValidAndUnique(t *testing.T) {
	caps := DefaultCapabilities()
	if len(caps) == 0 { t.Fatal("default capability registry must not be empty") }
	seen := map[string]bool{}
	for _, c := range caps {
		if !c.Valid() { t.Fatalf("invalid default capability: %#v", c) }
		if seen[c.ID] { t.Fatalf("duplicate capability: %q", c.ID) }
		seen[c.ID] = true
	}
}
