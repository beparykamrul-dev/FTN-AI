package security

// EvidenceRepository abstracts persistence for PostgreSQL, CockroachDB, or another backend.
type EvidenceRepository interface {
	Put(EvidenceRecord) error
	Get(id string) (EvidenceRecord, error)
}
