package data

// MigrationStep represents one controlled database migration/replication step.
type MigrationStep struct {
	From      string
	To        string
	Scope     string
	ReadOnly  bool
	Verified  bool
}

// MigrationPlan is intentionally declarative; an executor must validate and
// authorize it before changing production data.
type MigrationPlan struct {
	Steps []MigrationStep
}

func (p MigrationPlan) Valid() bool {
	for _, s := range p.Steps {
		if s.From == "" || s.To == "" || s.From == s.To {
			return false
		}
		if s.Scope != "global" && s.Scope != "regional" && s.Scope != "local" {
			return false
		}
	}
	return true
}
