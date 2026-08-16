package security

// PostgresEvidenceRepository is the DB-facing adapter boundary.
// The concrete SQL driver/connection pool is injected by the application layer.
type PostgresEvidenceRepository struct {
	PutFunc func(EvidenceRecord) error
	GetFunc func(string) (EvidenceRecord, error)
}

func (r PostgresEvidenceRepository) Put(e EvidenceRecord) error {
	if r.PutFunc == nil { return nil }
	return r.PutFunc(e)
}

func (r PostgresEvidenceRepository) Get(id string) (EvidenceRecord, error) {
	if r.GetFunc == nil { return EvidenceRecord{}, nil }
	return r.GetFunc(id)
}
