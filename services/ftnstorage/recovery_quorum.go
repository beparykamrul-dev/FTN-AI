package ftnstorage

// RecoveryQuorum describes the minimum independent verified sources needed
// before a recovery decision is considered sufficiently supported.
type RecoveryQuorum struct {
	Required int
}

func (q RecoveryQuorum) Satisfied(replicas []Replica) bool {
	if q.Required <= 0 {
		return false
	}
	count := 0
	for _, r := range replicas {
		if r.Healthy && r.Verified {
			count++
		}
	}
	return count >= q.Required
}
