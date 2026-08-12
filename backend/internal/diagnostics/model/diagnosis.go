package model

// Diagnosis contains evidence-backed diagnostic output.
type Diagnosis struct {
	IncidentID string   `json:"incident_id"`
	Cause      string   `json:"cause,omitempty"`
	Confidence float64  `json:"confidence"`
	Evidence   []string `json:"evidence,omitempty"`
	Impact     []string `json:"impact,omitempty"`
}
