package servicefoundry

import (
	"errors"
	"sort"
)

// ProviderEvidence is normalized evidence gathered about an ecosystem option.
// It intentionally stores observations, not copied provider implementation.
type ProviderEvidence struct {
	Provider       string   `json:"provider"`
	Features       []string `json:"features,omitempty"`
	Security       []string `json:"security,omitempty"`
	Architecture   []string `json:"architecture,omitempty"`
	Scalability    []string `json:"scalability,omitempty"`
	Observability  []string `json:"observability,omitempty"`
	Limitations    []string `json:"limitations,omitempty"`
	License        string   `json:"license,omitempty"`
	Score          float64  `json:"score"`
}

type AnalysisWeights struct {
	Security      float64
	Reliability   float64
	Performance   float64
	Maintainability float64
	Observability float64
}

// ProviderAnalysis produces a ranked, evidence-based comparison. Scores are
// supplied by trusted collectors/reviewers; this package does not pretend to
// discover facts that have not been provided.
type ProviderAnalysis struct {
	ServiceName string             `json:"service_name"`
	Weights     AnalysisWeights   `json:"weights"`
	Candidates  []ProviderEvidence `json:"candidates"`
	Best        string             `json:"best,omitempty"`
	Gaps        []string           `json:"gaps,omitempty"`
}

func RankProviders(a ProviderAnalysis) ProviderAnalysis {
	out := a
	out.Candidates = append([]ProviderEvidence(nil), a.Candidates...)
	sort.SliceStable(out.Candidates, func(i, j int) bool {
		return out.Candidates[i].Score > out.Candidates[j].Score
	})
	if len(out.Candidates) > 0 { out.Best = out.Candidates[0].Provider }
	return out
}

// SecurityGate is a release boundary. A service cannot be certified until all
// mandatory controls are explicitly satisfied.
type SecurityGate struct {
	ThreatModel      bool `json:"threat_model"`
	AuthN            bool `json:"authentication"`
	AuthZ            bool `json:"authorization"`
	InputValidation  bool `json:"input_validation"`
	SecretsIsolation bool `json:"secrets_isolation"`
	Audit            bool `json:"audit"`
	ResourceLimits   bool `json:"resource_limits"`
	StaticAnalysis   bool `json:"static_analysis"`
	Tests            bool `json:"tests"`
	DependencyAudit  bool `json:"dependency_audit"`
}

func (g SecurityGate) Validate() error {
	checks := map[string]bool{
		"threat model": g.ThreatModel, "authentication": g.AuthN,
		"authorization": g.AuthZ, "input validation": g.InputValidation,
		"secrets isolation": g.SecretsIsolation, "audit": g.Audit,
		"resource limits": g.ResourceLimits, "static analysis": g.StaticAnalysis,
		"tests": g.Tests, "dependency audit": g.DependencyAudit,
	}
	for name, ok := range checks {
		if !ok { return errors.New("security gate failed: " + name) }
	}
	return nil
}
