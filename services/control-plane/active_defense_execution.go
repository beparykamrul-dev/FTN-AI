package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type ActiveDefenseExecutionIntent struct {
	Action            WazuhActionClass `json:"action"`
	TargetAsset       string            `json:"target_asset"`
	TargetScope       string            `json:"target_scope"`
	DurationSeconds   int              `json:"duration_seconds"`
	Automatic         bool             `json:"automatic"`
	RequiresApproval  bool             `json:"requires_approval"`
	IdempotencyKey    string           `json:"idempotency_key"`
	SnapshotRequired  bool             `json:"snapshot_required"`
	VerificationRequired bool          `json:"verification_required"`
	Reason            string           `json:"reason"`
}

func BuildActiveDefenseExecutionIntent(alert WazuhAlert, assetRef string, owned bool) ActiveDefenseExecutionIntent {
	assetRef = strings.TrimSpace(assetRef)
	decision := BuildActiveDefenseDecision(alert, owned && assetRef != "")
	if decision.TargetScope != "ftn-owned-asset" {
		return ActiveDefenseExecutionIntent{
			Action: decision.Action, TargetScope: "none", RequiresApproval: true,
			Reason: decision.Reason, IdempotencyKey: ActiveDefenseIdempotencyKey(alert, "none"),
		}
	}
	intent := ActiveDefenseExecutionIntent{
		Action: decision.Action, TargetAsset: assetRef, TargetScope: decision.TargetScope,
		DurationSeconds: 900, Automatic: decision.Automatic, RequiresApproval: decision.RequiresApproval,
		SnapshotRequired: true, VerificationRequired: true, Reason: decision.Reason,
	}
	intent.IdempotencyKey = ActiveDefenseIdempotencyKey(alert, assetRef)
	return intent
}

func ActiveDefenseIdempotencyKey(alert WazuhAlert, assetRef string) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s", alert.ID, alert.RuleID, alert.Timestamp, assetRef, NormalizeWazuhSeverity(alert.Severity))
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
