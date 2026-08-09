package dns

import "fmt"

// PerlAdapter represents a compatibility boundary for existing Perl DNS tooling.
// The adapter describes capabilities; execution must be performed by an approved
// worker with explicit policy rather than arbitrary code from a request.
type PerlAdapter struct {
	Name    string
	Version string
}

func NewPerlAdapter(version string) *PerlAdapter {
	return &PerlAdapter{Name: "perl-dns", Version: version}
}

func (p *PerlAdapter) Validate() error {
	if p.Version == "" { return fmt.Errorf("perl version is required") }
	return nil
}
