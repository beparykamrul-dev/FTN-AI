package ai

import "github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"

// Recommendation is an advisory result. It is never an execution command.
type Recommendation struct {
	Summary    string   `json:"summary"`
	Action     string   `json:"action,omitempty"`
	Confidence float64  `json:"confidence"`
	Evidence   []string `json:"evidence,omitempty"`
	ApprovalRequired bool `json:"approval_required"`
}

type Advisor struct{}

// Advise converts a validated diagnosis into a conservative human-readable
// recommendation. It deliberately does not access infrastructure credentials.
func (Advisor) Advise(d model.Diagnosis) Recommendation {
	if d.Cause == "" {
		return Recommendation{
			Summary: "Insufficient evidence for a reliable root-cause recommendation.",
			Confidence: d.Confidence,
			Evidence: append([]string(nil), d.Evidence...),
			ApprovalRequired: true,
		}
	}
	return Recommendation{
		Summary: d.Cause,
		Action: "Review evidence and apply an approved remediation policy.",
		Confidence: d.Confidence,
		Evidence: append([]string(nil), d.Evidence...),
		ApprovalRequired: true,
	}
}
