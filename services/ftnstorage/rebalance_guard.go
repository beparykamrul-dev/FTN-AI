package ftnstorage

// RebalanceGuard keeps storage movement behind explicit health and approval
// requirements. It does not move bytes itself.
type RebalanceGuard struct {
	RequireApproval bool
	RequireHealthy  bool
}

func (g RebalanceGuard) Allowed(plan RebalancePlan, healthySource, healthyTarget, approved bool) bool {
	if !plan.Valid() || (g.RequireHealthy && (!healthySource || !healthyTarget)) {
		return false
	}
	if g.RequireApproval && !approved {
		return false
	}
	return true
}
