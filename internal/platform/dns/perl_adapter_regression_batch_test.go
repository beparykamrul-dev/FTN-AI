package dns

import "testing"

func TestPerlAdapterRejectsMissingVersion(t *testing.T) {
	if NewPerlAdapter("").Validate() == nil { t.Fatal("missing Perl version must be rejected") }
}

func TestPerlAdapterValidVersion(t *testing.T) {
	if NewPerlAdapter("5.40.0").Validate() != nil { t.Fatal("valid Perl version must be accepted") }
}
