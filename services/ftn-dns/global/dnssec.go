package global

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// DNSSECObservation is a provider-neutral result passed to the FTN DNS
// control plane after a DNS adapter has inspected a response chain.
type DNSSECObservation struct {
	Domain       string
	Validated    bool
	Authenticated bool
	Error        string
}

// DigestIdentity creates a stable, non-secret identity for an observed DNS
// object. It is useful for consistency comparisons without storing payloads.
func DigestIdentity(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

// DNSSECReady indicates whether an observation is safe to use for a
// validation-aware routing decision. The adapter remains responsible for
// performing the actual DNSSEC cryptographic validation.
func DNSSECReady(o DNSSECObservation) bool {
	return o.Error == "" && o.Validated && o.Authenticated
}
