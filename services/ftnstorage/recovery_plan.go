package ftnstorage

// RecoveryPlan is a declarative storage repair plan. Execution remains behind
// FTN's approval, audit, and verification boundaries.
type RecoveryPlan struct {
	NodeID       string
	Source       string
	Target       string
	SnapshotID   string
	BytesToCopy  uint64
	Verified     bool
}

func (p RecoveryPlan) Valid() bool {
	return p.NodeID != "" && p.Source != "" && p.Target != "" && p.SnapshotID != ""
}
