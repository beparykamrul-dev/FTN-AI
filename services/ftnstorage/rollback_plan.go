package ftnstorage

// RollbackPlan describes a move back to a verified snapshot without executing it.
type RollbackPlan struct {
	NodeID         string
	SnapshotID     string
	Reason         string
	RequireApproval bool
}

func (p RollbackPlan) Valid() bool {
	return p.NodeID != "" && p.SnapshotID != "" && p.RequireApproval
}
