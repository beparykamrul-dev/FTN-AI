package ftnstorage

// SnapshotRecoveryDecision chooses rollback only from a verified snapshot.
type SnapshotRecoveryDecision struct {
	NodeID     string
	SnapshotID string
	Verified   bool
	Reason     string
}

func (d SnapshotRecoveryDecision) Valid() bool {
	return d.NodeID != "" && d.SnapshotID != "" && d.Verified
}
