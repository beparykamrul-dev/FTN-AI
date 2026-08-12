package correlate

import (
	"sort"

	"github.com/beparykamrul-dev/FTN-AI/backend/internal/diagnostics/model"
)

// Correlator orders authorized evidence for deterministic incident analysis.
type Correlator struct{}

func (Correlator) Collect(incident model.Incident, evidence []model.Evidence) []model.Evidence {
	allowed := make(map[string]struct{}, len(incident.EvidenceIDs))
	for _, id := range incident.EvidenceIDs {
		allowed[id] = struct{}{}
	}

	out := make([]model.Evidence, 0, len(evidence))
	for _, item := range evidence {
		if _, ok := allowed[item.ID]; !ok {
			continue
		}
		out = append(out, item)
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Timestamp < out[j].Timestamp
	})
	return out
}
