package data

// TransactionGuard expresses safety checks required before a production data
// migration or replication operation is executed.
type TransactionGuard struct {
	PlanVersion uint64
	Approved    bool
	Verified    bool
}

func (g TransactionGuard) CanExecute(version uint64) bool {
	return g.Approved && g.Verified && g.PlanVersion == version
}
