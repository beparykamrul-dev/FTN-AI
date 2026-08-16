package observability

// MultipathPlan describes bounded parallel use of independent healthy paths.
type MultipathPlan struct {
	PathIDs []string
	MaxParallel uint32
}

func (p MultipathPlan) Valid() bool {
	return len(p.PathIDs) > 0 && p.MaxParallel > 0
}
