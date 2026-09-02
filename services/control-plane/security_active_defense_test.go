package main

import "testing"

func TestActiveDefenseNeverTargetsUnverifiedExternalAssets(t *testing.T) {
	d := BuildActiveDefenseDecision(WazuhAlert{Severity: "critical"}, false)
	if d.TargetScope != "none" || d.Automatic || !d.RequiresApproval {
		t.Fatalf("unsafe target decision: %#v", d)
	}
}

func TestActiveDefenseContainsCriticalFTNAsset(t *testing.T) {
	d := BuildActiveDefenseDecision(WazuhAlert{Severity: "critical"}, true)
	if d.Action != WazuhHealthRecover || !d.Automatic || d.TargetScope != "ftn-owned-asset" {
		t.Fatalf("unexpected decision: %#v", d)
	}
}

func TestActiveDefenseMediumIsBounded(t *testing.T) {
	d := BuildActiveDefenseDecision(WazuhAlert{Severity: "medium"}, true)
	if d.Action != WazuhAlertAction || !d.Automatic {
		t.Fatalf("unexpected decision: %#v", d)
	}
}
