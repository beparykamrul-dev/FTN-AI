package ftnstorage

// HealthPolicy defines bounded automatic storage-health actions.
type HealthPolicy struct {
	FailureThreshold uint32
	QuarantineAfter  uint32
	RequireVerified  bool
}

func (p HealthPolicy) Valid() bool {
	return p.FailureThreshold > 0 && p.QuarantineAfter >= p.FailureThreshold
}

func (p HealthPolicy) ShouldQuarantine(failures uint32) bool {
	return p.Valid() && failures >= p.QuarantineAfter
}
