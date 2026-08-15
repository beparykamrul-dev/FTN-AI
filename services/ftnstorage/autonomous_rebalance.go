package ftnstorage

// RebalanceDecision is the safe, declarative result of storage placement evaluation.
type RebalanceDecision struct {
	Plan        RebalancePlan
	HealthRisk  float64
	Allowed     bool
	Reason      string
}

func BuildRebalanceDecision(plan RebalancePlan, risk float64) RebalanceDecision {
	allowed := plan.Valid() && risk < 0.7
	return RebalanceDecision{Plan: plan, HealthRisk: risk, Allowed: allowed, Reason: "health-aware-placement-evaluation"}
}
