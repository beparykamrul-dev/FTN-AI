package fiber

import "time"

type FiberNode struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status string `json:"status"`
}

type FiberLink struct {
	ID string `json:"id"`
	From string `json:"from"`
	To string `json:"to"`
	DistanceMeters float64 `json:"distance_meters"`
	Status string `json:"status"`
}

type RecoveryCandidate struct {
	LinkID string `json:"link_id"`
	Confidence float64 `json:"confidence"`
	Reason string `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyzeFailure creates a recovery recommendation from observed fiber state.
// It does not execute physical or network changes; execution remains behind
// the FTN approval workflow.
func AnalyzeFailure(link FiberLink, evidence []string) RecoveryCandidate {
	confidence := 0.50
	if link.Status == "cut" || link.Status == "down" { confidence += 0.25 }
	if len(evidence) > 0 { confidence += 0.05 * float64(len(evidence)) }
	if confidence > 0.95 { confidence = 0.95 }
	return RecoveryCandidate{LinkID: link.ID, Confidence: confidence, Reason: "fiber failure evidence requires approved recovery", CreatedAt: time.Now().UTC()}
}
