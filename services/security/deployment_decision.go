package security

// DeploymentDecision is the final policy result before deployment execution.
type DeploymentDecision struct {
	Allowed bool
	Reason  string
}

func DecideDeployment(risk RiskAggregate, approval ApprovalRecord, nowUnix int64) DeploymentDecision {
	if risk.Critical > 0 || risk.High > 0 {
		if !approval.Valid(nowUnix) {
			return DeploymentDecision{Allowed: false, Reason: "blocking-findings-require-valid-approval"}
		}
	}
	return DeploymentDecision{Allowed: true, Reason: "security-policy-satisfied"}
}
