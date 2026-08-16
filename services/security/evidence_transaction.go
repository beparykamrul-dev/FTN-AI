package security

// EvidenceTransaction describes the atomic persistence boundary used by a DB adapter.
type EvidenceTransaction struct {
	Isolation string
	Retryable bool
}

func (t EvidenceTransaction) Valid() bool {
	return t.Isolation != ""
}
