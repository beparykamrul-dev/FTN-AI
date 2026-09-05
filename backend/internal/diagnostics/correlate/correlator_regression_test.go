package correlate

import (
	"testing"
	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

func TestCorrelatorFiltersUnauthorizedAndDuplicateEvidence(t *testing.T) {
	incident := model.Incident{ID:"i1", Severity:"high", EvidenceIDs:[]string{"e1","e1"}}
	input := []model.Evidence{{ID:"e1",Kind:"probe",Source:"n1"},{ID:"e1",Kind:"probe",Source:"n1"},{ID:"e2",Kind:"probe",Source:"n1"}}
	out := (Correlator{}).Collect(incident, input)
	if len(out) != 1 || out[0].ID != "e1" { t.Fatalf("unexpected correlation output: %#v", out) }
}
