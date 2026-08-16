package security

// RouteEvidence selects the first valid enabled backend deterministically.
func RouteEvidence(backends []EvidenceBackend) (EvidenceBackend, bool) {
	for _, b := range backends {
		if b.Valid() { return b, true }
	}
	return EvidenceBackend{}, false
}
