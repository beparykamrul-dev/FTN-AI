package security

// ArtifactPolicy defines release-artifact requirements before deployment.
type ArtifactPolicy struct {
	RequireAttestation bool
	RequireFingerprint bool
}

func (p ArtifactPolicy) Allows(attested, fingerprinted bool) bool {
	if p.RequireAttestation && !attested { return false }
	if p.RequireFingerprint && !fingerprinted { return false }
	return true
}
