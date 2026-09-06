package dns

import "testing"

func TestPerlAdapterValidation(t *testing.T) {
 if err := NewPerlAdapter("5.38").Validate(); err != nil { t.Fatalf("valid Perl adapter rejected: %v", err) }
 if err := (*PerlAdapter)(nil).Validate(); err == nil { t.Fatal("nil Perl adapter accepted") }
 if err := NewPerlAdapter("5.38\n").Validate(); err == nil { t.Fatal("control character in Perl version accepted") }
}
