package security

// EvidenceBackend identifies the centralized persistence backend selected by FTN.
type EvidenceBackend struct {
	Name string
	Role string
	Enabled bool
}

func (b EvidenceBackend) Valid() bool {
	return b.Name != "" && b.Role != "" && b.Enabled
}
