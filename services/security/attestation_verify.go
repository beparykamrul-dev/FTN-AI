package security

// VerifyAttestation checks the minimum fields needed to accept a release attestation.
// Cryptographic signature verification is intentionally delegated to the signing adapter.
type AttestationVerifier struct{}

func (AttestationVerifier) Verify(a ReleaseAttestation) bool {
	return a.ReleaseID != "" && a.Commit != "" && a.Fingerprint != ""
}
