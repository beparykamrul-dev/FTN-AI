package ftnstorage

// GCPlan is a declarative reclaim plan. It does not delete storage itself.
type GCPlan struct {
	PolicyVersion uint64
	ChunkIDs     []string
}

func (p GCPlan) Valid() bool {
	return p.PolicyVersion > 0 && len(p.ChunkIDs) > 0
}
