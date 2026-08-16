package security

// GateStatus is the normalized state of the FTN security gate.
type GateStatus struct {
	Allowed bool
	Score int
	Reason string
}

func BuildGateStatus(findings []Finding, policy RiskPolicy, approved bool) GateStatus {
	risk := AggregateRisk(findings)
	score := SecurityScore(findings)
	if !policy.Allows(risk, approved) {
		return GateStatus{Allowed: false, Score: score, Reason: "risk-policy-blocked"}
	}
	return GateStatus{Allowed: true, Score: score, Reason: "risk-policy-passed"}
}
