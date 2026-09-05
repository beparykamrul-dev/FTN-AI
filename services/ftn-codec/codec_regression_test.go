package codec

import "testing"

func TestCapabilityValidationRejectsMissingIdentity(t *testing.T) {
	if (Capability{Class:"network"}).Valid() { t.Fatal("capability without id must be invalid") }
}

func TestCapabilityNormalizationTrimsIdentity(t *testing.T) {
	c := (Capability{ID:" id ", Class:" network "}).Normalize()
	if c.ID != "id" || c.Class != "network" { t.Fatalf("unexpected normalization: %#v", c) }
}
