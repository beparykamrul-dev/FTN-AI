package ftnwan

// FailoverPlan describes a non-destructive switch to the next healthy path.
type FailoverPlan struct {
	Current string
	Next    string
	Reason  string
}

// BuildFailover selects the best healthy path other than the current path.
func BuildFailover(paths []Path, current string) (FailoverPlan, bool) {
	best, ok := SelectPath(paths)
	if !ok || best.ID == current {
		return FailoverPlan{}, false
	}
	return FailoverPlan{Current: current, Next: best.ID, Reason: "health or path quality changed"}, true
}
