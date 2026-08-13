package adapters

import (
	"crypto/sha256"
	"fmt"
)

// VerifyDSDigest verifies the digest binding between a canonical DNSKEY value
// and a DS record. The caller must supply canonical DNSSEC wire material.
// Supported digest type 2 is SHA-256.
func VerifyDSDigest(canonicalDNSKEY []byte, expected []byte, digestType uint8) error {
	if len(canonicalDNSKEY) == 0 || len(expected) == 0 {
		return fmt.Errorf("empty DNSSEC digest input")
	}
	if digestType != 2 {
		return fmt.Errorf("unsupported DS digest type: %d", digestType)
	}
	digest := sha256.Sum256(canonicalDNSKEY)
	if len(expected) != len(digest) {
		return fmt.Errorf("DS digest length mismatch")
	}
	for i := range digest {
		if digest[i] != expected[i] {
			return fmt.Errorf("DS digest mismatch")
		}
	}
	return nil
}

// VerifyChainBinding verifies the DS-to-DNSKEY binding. Signature verification
// over RRSIG data remains a separate operation because it requires canonical
// RRset construction and algorithm-specific public-key handling.
func VerifyChainBinding(dnskey []byte, dsDigest []byte, digestType uint8) bool {
	return VerifyDSDigest(dnskey, dsDigest, digestType) == nil
}
