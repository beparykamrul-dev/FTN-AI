package ftnwan

// FailoverPlan describes a non-destructive switch to the next healthy path.
type FailoverPlan struct {
	Current string
	Next    string
	Reason  string
}

func BuildFailover(paths []Path, current string) (FailoverPlan, bool) {
	for _, p := range SelectBest(paths) {
		if p.ID != current {
			return FailoverPlan{Current: current, Next: p.ID, Reason: "current path is not preferred"}, true
		}
	}
	return FailoverPlan{}, false
}
