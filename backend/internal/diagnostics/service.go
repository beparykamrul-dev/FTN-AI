package diagnostics

import "github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"

// Service coordinates evidence-backed diagnosis. It has no privileged execution path.
type Service struct{}

// Diagnose returns a conservative baseline diagnosis. Production adapters will supply
// telemetry correlation and dependency graph evidence without exposing secrets.
func (s Service) Diagnose(i model.Incident) model.Diagnosis {
	return model.Diagnosis{
		IncidentID: i.ID,
		Confidence: 0,
		Evidence:   append([]string(nil), i.EvidenceIDs...),
	}
}
