package security

// ReleaseChainEntry provides a compact evidence-chain link for a release.
type ReleaseChainEntry struct {
	ReleaseID string
	ParentFingerprint string
	CurrentFingerprint string
}

func (e ReleaseChainEntry) Valid() bool {
	return e.ReleaseID != "" && e.CurrentFingerprint != ""
}
