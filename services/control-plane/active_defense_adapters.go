package main

import "strings"

// ActiveDefenseAdapter is an execution boundary for FTN-owned containment.
// Implementations must perform ownership checks again at execution time and
// must not accept arbitrary shell commands or external targets.
type ActiveDefenseAdapter string

const (
	AdapterNFTables ActiveDefenseAdapter = "nftables"
	AdapterXDP      ActiveDefenseAdapter = "xdp"
	AdapterMikroTik ActiveDefenseAdapter = "mikrotik"
)

type ActiveDefenseActionPlan struct {
	Adapter            ActiveDefenseAdapter `json:"adapter"`
	Operation          string               `json:"operation"`
	TargetAsset        string               `json:"target_asset"`
	DurationSeconds    int                  `json:"duration_seconds"`
	SnapshotRequired   bool                 `json:"snapshot_required"`
	VerificationNeeded bool                 `json:"verification_needed"`
	Automatic          bool                 `json:"automatic"`
	RequiresApproval   bool                 `json:"requires_approval"`
}

func BuildActiveDefenseActionPlan(intent ActiveDefenseExecutionIntent) (ActiveDefenseActionPlan, bool) {
	if intent.TargetScope != "ftn-owned-asset" || strings.TrimSpace(intent.TargetAsset) == "" {
		return ActiveDefenseActionPlan{}, false
	}
	if intent.DurationSeconds <= 0 || intent.DurationSeconds > 3600 {
		return ActiveDefenseActionPlan{}, false
	}
	// Active-defense mutations are privileged even when an upstream decision
	// was marked automatic; the execution boundary must remain approval-gated.
	plan := ActiveDefenseActionPlan{
		TargetAsset: intent.TargetAsset,
		DurationSeconds: intent.DurationSeconds,
		SnapshotRequired: true,
		VerificationNeeded: true,
		Automatic: intent.Automatic,
		RequiresApproval: true,
	}
	switch intent.Action {
	case WazuhHealthRecover:
		plan.Adapter, plan.Operation = AdapterNFTables, "temporary-containment"
	case WazuhAlertAction:
		plan.Adapter, plan.Operation = AdapterNFTables, "rate-limit"
	default:
		return ActiveDefenseActionPlan{}, false
	}
	return plan, true
}
