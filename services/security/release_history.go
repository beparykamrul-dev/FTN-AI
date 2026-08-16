package security

// ReleaseHistory records immutable release verification outcomes.
type ReleaseHistory struct {
	ReleaseID string
	Commit string
	Policy PolicyVersion
	Score int
	Allowed bool
}

func (h ReleaseHistory) Valid() bool {
	return h.ReleaseID != "" && h.Commit != "" && h.Policy.Valid()
}
