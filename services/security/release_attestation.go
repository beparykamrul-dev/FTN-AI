package security

// ReleaseAttestation binds a release to the security scan and policy state
// that authorized it. It contains no source code or secret material.
type ReleaseAttestation struct {
	ReleaseID string
	Commit string
	ScanRunID string
	Policy PolicyVersion
	Fingerprint string
}

func (a ReleaseAttestation) Valid() bool {
	return a.ReleaseID != "" && a.Commit != "" && a.ScanRunID != "" && a.Policy.Valid() && a.Fingerprint != ""
}
