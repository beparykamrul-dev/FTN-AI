package security

// ReleaseRejection is a structured, auditable block reason.
type ReleaseRejection struct {
	ReleaseID string
	Reason string
	Policy PolicyVersion
}

func (r ReleaseRejection) Valid() bool {
	return r.ReleaseID != "" && r.Reason != "" && r.Policy.Valid()
}
