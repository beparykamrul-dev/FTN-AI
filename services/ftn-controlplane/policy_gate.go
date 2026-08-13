package controlplane

// PolicyGate is the mandatory boundary between analysis and deployment.
type PolicyGate struct {
	RequireHealthy bool
	RequireKnownServer bool
	RequireExplicitApproval bool
}

// DeploymentPlan is an immutable-by-convention plan produced after policy
// validation. Execution is intentionally outside this layer.
type DeploymentPlan struct {
	ServerID string
	Services []string
	Approved bool
	Reason string
}

func (g PolicyGate) Validate(a AnalysisResult, serverKnown, approved bool) DeploymentPlan {
	if g.RequireKnownServer && !serverKnown {
		return DeploymentPlan{ServerID: a.ServerID, Reason: "server is not enrolled"}
	}
	if g.RequireHealthy && !a.Healthy {
		return DeploymentPlan{ServerID: a.ServerID, Reason: "server health policy rejected deployment"}
	}
	if g.RequireExplicitApproval && !approved {
		return DeploymentPlan{ServerID: a.ServerID, Reason: "explicit approval required"}
	}
	return DeploymentPlan{ServerID: a.ServerID, Services: append([]string(nil), a.Desired.Services...), Approved: true, Reason: "deployment policy accepted"}
}
