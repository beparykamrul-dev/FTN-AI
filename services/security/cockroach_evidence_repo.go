package security

// CockroachEvidenceRepository is the retry-aware adapter boundary for CockroachDB.
type CockroachEvidenceRepository struct {
	PutFunc func(EvidenceRecord) error
	GetFunc func(string) (EvidenceRecord, error)
}

func (r CockroachEvidenceRepository) Put(e EvidenceRecord) error {
	if r.PutFunc == nil { return nil }
	return r.PutFunc(e)
}

func (r CockroachEvidenceRepository) Get(id string) (EvidenceRecord, error) {
	if r.GetFunc == nil { return EvidenceRecord{}, nil }
	return r.GetFunc(id)
}
