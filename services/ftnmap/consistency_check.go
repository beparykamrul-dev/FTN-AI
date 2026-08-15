package ftnmap

// ConsistencyReport describes deterministic topology consistency checks.
type ConsistencyReport struct {
	ValidEvents uint64
	InvalidEvents uint64
	Consistent   bool
}

func CheckConsistency(events []TopologyEvent) ConsistencyReport {
	var r ConsistencyReport
	for _, e := range events {
		if e.Valid() { r.ValidEvents++ } else { r.InvalidEvents++ }
	}
	r.Consistent = r.InvalidEvents == 0
	return r
}
