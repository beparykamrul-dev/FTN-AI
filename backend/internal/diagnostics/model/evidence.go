package model

import (
	"math"
	"strings"
)

// Evidence is a normalized observation supporting or weakening a diagnosis.
type Evidence struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	Source string `json:"source"`
	Timestamp string `json:"timestamp"`
	Value float64 `json:"value,omitempty"`
	Summary string `json:"summary,omitempty"`
	Redacted bool `json:"redacted"`
}

func (e Evidence) Valid() bool {
	return strings.TrimSpace(e.ID) != "" && strings.TrimSpace(e.Kind) != "" && strings.TrimSpace(e.Source) != "" && !math.IsNaN(e.Value) && !math.IsInf(e.Value, 0)
}

func (e Evidence) Normalize() Evidence {
	e.ID = strings.TrimSpace(e.ID); e.Kind = strings.TrimSpace(e.Kind); e.Source = strings.TrimSpace(e.Source); e.Timestamp = strings.TrimSpace(e.Timestamp); e.Summary = strings.TrimSpace(e.Summary)
	return e
}
