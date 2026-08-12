package rootcause

import (
	"sort"

	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

// Candidate is a conservative root-cause hypothesis. It is not an execution command.
type Candidate struct {
	Cause      string
	Confidence float64
	Evidence   []string
	Impact     []string
}

type Analyzer struct{}

// Find ranks simple evidence-backed candidates. Production adapters can add
// domain-specific rules without granting the analyzer privileged execution.
func (Analyzer) Find(evidence []model.Evidence, dependencies []model.Dependency) []Candidate {
	candidates := make([]Candidate, 0)
	for _, item := range evidence {
		if item.Summary == "" {
			continue
		}
		candidates = append(candidates, Candidate{
			Cause:      item.Summary,
			Confidence: 0.25,
			Evidence:   []string{item.ID},
		})
	}

	for _, dep := range dependencies {
		if dep.Critical && !dep.Healthy {
			candidates = append(candidates, Candidate{
				Cause:      "critical dependency degraded",
				Confidence: 0.5,
				Impact:     []string{dep.From, dep.To},
			})
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Confidence > candidates[j].Confidence
	})
	return candidates
}

func (Analyzer) Validate(candidates []Candidate, evidence []model.Evidence) model.Diagnosis {
	if len(candidates) == 0 {
		return model.Diagnosis{Confidence: 0}
	}
	best := candidates[0]
	return model.Diagnosis{
		Cause:      best.Cause,
		Confidence: best.Confidence,
		Evidence:   append([]string(nil), best.Evidence...),
		Impact:     append([]string(nil), best.Impact...),
	}
}
