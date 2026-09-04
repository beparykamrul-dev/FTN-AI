package gis

import "time"

type RecoveryStep struct {
	ID string `json:"id"`
	Action string `json:"action"`
	Target string `json:"target"`
	Reason string `json:"reason"`
	RequiresApproval bool `json:"requires_approval"`
}

type RecoveryPlan struct {
	ID string `json:"id,omitempty"`
	RootAsset string `json:"root_asset,omitempty"`
	FailedAssetID string `json:"failed_asset_id,omitempty"`
	Risk string `json:"risk,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Affected []ImpactNode `json:"affected,omitempty"`
	Candidates []RecoveryCandidate `json:"candidates,omitempty"`
	Steps []RecoveryStep `json:"steps,omitempty"`
	RequiresApproval bool `json:"requires_approval"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// BuildRecoveryPlan produces an auditable plan only. It does not execute any
// network or physical remediation action.
func BuildRecoveryPlan(root string, impact ImpactResult, confidence float64) RecoveryPlan {
	if confidence < 0 { confidence = 0 }; if confidence > 1 { confidence = 1 }
	steps := make([]RecoveryStep, 0, len(impact.Affected))
	for _, n := range impact.Affected {
		steps = append(steps, RecoveryStep{ID:n.AssetID, Action:"verify_and_restore", Target:n.AssetID, Reason:"affected_by_root_failure", RequiresApproval:true})
	}
	risk := "medium"
	if len(impact.Affected) >= 50 { risk = "high" }
	now := time.Now().UTC()
	return RecoveryPlan{ID:root+"-"+now.Format("20060102150405"), RootAsset:root, Risk:risk, Confidence:confidence, Affected:impact.Affected, Steps:steps, RequiresApproval:true, CreatedAt:now}
}
