package adapters

import (
	"crypto"
	"crypto/sha256"
	"fmt"
)

// RRSIGVerifier is the crypto boundary for DNSSEC signature validation.
// Implementations must construct the canonical RRset and DNSSEC signed data
// according to the DNSSEC specification before calling Verify.
type RRSIGVerifier interface {
	Verify(algorithm uint8, publicKey, signedData, signature []byte) error
}

// HashSignedData provides a deterministic digest for diagnostics and cache
// identity. It is not a substitute for DNSSEC signature verification.
func HashSignedData(data []byte) [32]byte { return sha256.Sum256(data) }

// VerifyRRSIG delegates algorithm-specific verification to a trusted crypto
// implementation. An absent verifier is always treated as failure.
func VerifyRRSIG(v RRSIGVerifier, algorithm uint8, publicKey, signedData, signature []byte) error {
	if v == nil {
		return fmt.Errorf("DNSSEC RRSIG verifier is not configured")
	}
	if len(publicKey) == 0 || len(signedData) == 0 || len(signature) == 0 {
		return fmt.Errorf("incomplete RRSIG verification input")
	}
	if algorithm == 0 {
		return fmt.Errorf("invalid DNSSEC algorithm")
	}
	return v.Verify(algorithm, publicKey, signedData, signature)
}

// CryptoHash exposes the standard crypto hash type for concrete adapters that
// need to map DNSSEC algorithms to Go's crypto primitives.
func CryptoHash(h crypto.Hash) crypto.Hash { return h }
