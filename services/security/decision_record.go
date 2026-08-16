package security

// DecisionRecord makes a deployment security decision reproducible and auditable.
type DecisionRecord struct {
	RunID string
	Policy PolicyVersion
	Allowed bool
	Score int
	Reason string
}

func (r DecisionRecord) Valid() bool {
	return r.RunID != "" && r.Policy.Valid() && r.Reason != ""
}
