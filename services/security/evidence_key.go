package security

// EvidenceKey defines the stable lookup identity for persisted security evidence.
type EvidenceKey struct {
	Kind string
	ID string
}

func (k EvidenceKey) Valid() bool { return k.Kind != "" && k.ID != "" }
