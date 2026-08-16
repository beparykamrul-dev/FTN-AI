package security

// EvidenceRecord is the DB-neutral persisted representation of security evidence.
type EvidenceRecord struct {
	ID string
	Kind string
	RunID string
	ReleaseID string
	Fingerprint string
	PayloadRef string
	CreatedUnix int64
}

func (r EvidenceRecord) Valid() bool {
	return r.ID != "" && r.Kind != "" && r.RunID != "" && r.CreatedUnix > 0
}
