package adapters

import "fmt"

// RRSIGPolicy defines the minimum requirements for accepting a signed RRset.
type RRSIGPolicy struct {
	RequireSignature bool
	AllowedAlgorithms map[uint8]bool
}

// RRSIGResult is the normalized result returned by a cryptographic verifier.
type RRSIGResult struct {
	Present   bool
	Verified  bool
	Algorithm uint8
	Error     string
}

// EvaluateRRSIG applies the policy without performing cryptographic work.
// Verification must already have been performed by a trusted crypto adapter.
func EvaluateRRSIG(policy RRSIGPolicy, result RRSIGResult) error {
	if policy.RequireSignature && !result.Present {
		return fmt.Errorf("DNSSEC signature required")
	}
	if !result.Present {
		return nil
	}
	if !result.Verified {
		if result.Error != "" { return fmt.Errorf("RRSIG verification failed: %s", result.Error) }
		return fmt.Errorf("RRSIG verification failed")
	}
	if len(policy.AllowedAlgorithms) > 0 && !policy.AllowedAlgorithms[result.Algorithm] {
		return fmt.Errorf("RRSIG algorithm %d rejected by policy", result.Algorithm)
	}
	return nil
}
