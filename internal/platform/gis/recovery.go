package gis

import "time"

type RecoveryAction struct {
	ID string `json:"id"`
	AssetID string `json:"asset_id"`
	Kind string `json:"kind"`
	Reason string `json:"reason"`
	Confidence float64 `json:"confidence"`
	RequiresApproval bool `json:"requires_approval"`
	CreatedAt time.Time `json:"created_at"`
}

type RecoveryEngine struct{}

func NewRecoveryEngine() *RecoveryEngine { return &RecoveryEngine{} }

// Propose creates an auditable recovery proposal. It does not modify network
// state. Any physical repair or service-affecting operation requires explicit
// approval and a downstream authorized executor.
func (e *RecoveryEngine) Propose(asset FiberAsset, reason string, confidence float64) RecoveryAction {
	if confidence < 0 { confidence = 0 }; if confidence > 1 { confidence = 1 }
	return RecoveryAction{AssetID: asset.ID, Kind: string(asset.Type), Reason: reason, Confidence: confidence, RequiresApproval: true, CreatedAt: time.Now().UTC()}
}
