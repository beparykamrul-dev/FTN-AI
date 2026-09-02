package main

import "testing"

func TestActiveDefenseActionPlanRejectsUnownedTarget(t *testing.T) {
	intent := ActiveDefenseExecutionIntent{Action: WazuhHealthRecover, TargetAsset: "attacker", TargetScope: "none", DurationSeconds: 900}
	if _, ok := BuildActiveDefenseActionPlan(intent); ok {
		t.Fatal("unowned target must never produce an execution plan")
	}
}

func TestActiveDefenseActionPlanUsesBoundedAdapters(t *testing.T) {
	intent := ActiveDefenseExecutionIntent{Action: WazuhHealthRecover, TargetAsset: "ftn-router-1", TargetScope: "ftn-owned-asset", DurationSeconds: 900, Automatic: true}
	plan, ok := BuildActiveDefenseActionPlan(intent)
	if !ok || plan.Adapter != AdapterNFTables || plan.Operation != "temporary-containment" || !plan.SnapshotRequired || !plan.VerificationNeeded {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}

func TestActiveDefenseActionPlanRejectsUnsafeDuration(t *testing.T) {
	intent := ActiveDefenseExecutionIntent{Action: WazuhHealthRecover, TargetAsset: "ftn-router-1", TargetScope: "ftn-owned-asset", DurationSeconds: 3601}
	if _, ok := BuildActiveDefenseActionPlan(intent); ok {
		t.Fatal("duration above safety limit must be rejected")
	}
}
