package model

import "math"

// Diagnosis contains evidence-backed diagnostic output.
type Diagnosis struct { IncidentID string `json:"incident_id"`; Cause string `json:"cause,omitempty"`; Confidence float64 `json:"confidence"`; Evidence []string `json:"evidence,omitempty"`; Impact []string `json:"impact,omitempty"` }

func (d Diagnosis) Valid() bool { return d.IncidentID != "" && !math.IsNaN(d.Confidence) && !math.IsInf(d.Confidence, 0) && d.Confidence >= 0 && d.Confidence <= 1 }
