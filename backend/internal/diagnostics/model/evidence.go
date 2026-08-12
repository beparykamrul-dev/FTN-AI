package model

// Evidence is a normalized observation supporting or weakening a diagnosis.
type Evidence struct {
	ID        string  `json:"id"`
	Kind      string  `json:"kind"`
	Source    string  `json:"source"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value,omitempty"`
	Summary   string  `json:"summary,omitempty"`
	Redacted  bool    `json:"redacted"`
}
