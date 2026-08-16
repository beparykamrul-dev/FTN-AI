package security

// ReleaseVerification is the normalized final pre-deployment security state.
type ReleaseVerification struct {
	ReleaseID string
	Attested bool
	Fingerprint string
	PolicyOK bool
}

func (r ReleaseVerification) Ready() bool {
	return r.ReleaseID != "" && r.Attested && r.Fingerprint != "" && r.PolicyOK
}
