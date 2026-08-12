package servicefoundry

import "sort"

// ArchitecturePlan converts analyzed capabilities into FTN-owned architectural
// boundaries. It never copies provider implementation; provider evidence only
// influences requirements and priorities.
type ArchitecturePlan struct {
	ServiceName string   `json:"service_name"`
	Boundaries  []string `json:"boundaries"`
	Capabilities []string `json:"capabilities"`
	SecurityControls []string `json:"security_controls"`
	OperationalControls []string `json:"operational_controls"`
	ExternalHostedRequired bool `json:"external_hosted_required"`
}

func GenerateNativeArchitecture(a ProviderAnalysis) ArchitecturePlan {
	plan := ArchitecturePlan{
		ServiceName: a.ServiceName,
		Boundaries: []string{"api", "control", "data", "events", "observability"},
		SecurityControls: []string{"authentication", "authorization", "input-validation", "audit", "resource-limits", "secrets-isolation"},
		OperationalControls: []string{"health", "metrics", "logs", "graceful-shutdown", "recovery", "data-consistency"},
		ExternalHostedRequired: false,
	}
	seen := map[string]bool{}
	for _, c := range a.Candidates {
		for _, feature := range c.Features {
			if feature != "" && !seen[feature] { plan.Capabilities = append(plan.Capabilities, feature); seen[feature] = true }
		}
	}
	sort.Strings(plan.Capabilities)
	return plan
}
