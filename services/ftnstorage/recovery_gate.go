package ftnstorage

// RecoveryGate prevents storage mutation unless the recovery decision is
// verified and explicitly authorized by the caller's control plane.
type RecoveryGate struct {
	ApprovedPlanVersion uint64
	Verified            bool
	Approved            bool
}

func (g RecoveryGate) Open(planVersion uint64) bool {
	return g.Verified && g.Approved && planVersion != 0 && planVersion == g.ApprovedPlanVersion
}
