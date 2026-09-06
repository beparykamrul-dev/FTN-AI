package codec

import "testing"

func TestCapabilityValidationBoundsIdentity(t *testing.T) {
	if (Capability{ID: "", Class: "transport"}).Valid() {
		t.Fatal("empty capability id must be rejected")
	}
	if !(Capability{ID: "binary-framing", Class: "transport"}).Valid() {
		t.Fatal("valid capability was rejected")
	}
}

func TestJobValidationRequiresInput(t *testing.T) {
	if err := (Job{CapabilityID: "binary-framing"}).Valid(); err == nil {
		t.Fatal("missing input URI must be rejected")
	}
}
