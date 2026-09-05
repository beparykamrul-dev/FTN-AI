package ai

import (
	"math"
	"testing"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

func TestAdvisorClampsInvalidConfidence(t *testing.T) {
	r := (Advisor{}).Advise(model.Diagnosis{Confidence:math.NaN()})
	if r.Confidence != 0 { t.Fatalf("expected confidence clamp, got %v", r.Confidence) }
	if !r.ApprovalRequired { t.Fatal("diagnostic remediation must remain approval gated") }
}

func TestAdvisorBoundsEvidence(t *testing.T) {
	evidence := []string{" a ", "a", "", string(make([]byte, 2050))}
	r := (Advisor{}).Advise(model.Diagnosis{Cause:"cause", Evidence:evidence, Confidence:.8})
	if len(r.Evidence) != 1 || r.Evidence[0] != "a" { t.Fatalf("unexpected evidence filtering: %#v", r.Evidence) }
}
