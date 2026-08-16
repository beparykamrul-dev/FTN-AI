package security

// GateDecision is the policy result used by FTN deployment controls.
type GateDecision struct {
	Allowed bool
	Reason  string
}

// Evaluate returns a conservative decision: invalid findings are rejected;
// valid findings at high/critical severity block deployment.
func Evaluate(findings []Finding) GateDecision {
	for _, f := range findings {
		if !f.Valid() { return GateDecision{Allowed: false, Reason: "invalid-security-finding"} }
		if f.Severity == "critical" || f.Severity == "high" {
			return GateDecision{Allowed: false, Reason: "blocking-security-finding"}
		}
	}
	return GateDecision{Allowed: true, Reason: "security-gate-passed"}
}
