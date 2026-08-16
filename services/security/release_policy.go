package security

// ReleasePolicy defines minimum security requirements for a deployable release.
type ReleasePolicy struct {
	RequireAttestation bool
	RequireFingerprint bool
	MinSecurityScore int
}

func (p ReleasePolicy) Allows(v ReleaseVerification, score int) bool {
	if p.RequireAttestation && !v.Attested { return false }
	if p.RequireFingerprint && v.Fingerprint == "" { return false }
	return score >= p.MinSecurityScore
}
