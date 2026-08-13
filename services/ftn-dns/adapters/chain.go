package adapters

// TrustChainStep describes one link in DNSSEC chain evaluation.
type TrustChainStep struct {
	Domain      string
	DSPresent   bool
	DNSKEYPresent bool
	RRSIGPresent bool
	Validated   bool
	Reason      string
}

// ValidateTrustChain performs a conservative structural chain check. It does
// not claim cryptographic validity; signature verification must be supplied by
// a DNSSEC crypto adapter before Validated can be trusted for routing.
func ValidateTrustChain(steps []TrustChainStep) bool {
	if len(steps) == 0 {
		return false
	}
	for _, step := range steps {
		if !step.DSPresent || !step.DNSKEYPresent || !step.RRSIGPresent || !step.Validated {
			return false
		}
	}
	return true
}
