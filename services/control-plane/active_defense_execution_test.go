package main

import "testing"

func TestActiveDefenseExecutionRequiresVerifiedOwnedAsset(t *testing.T) {
	intent := BuildActiveDefenseExecutionIntent(WazuhAlert{ID: "a1", RuleID: "r1", Severity: "critical", Timestamp: "2026-09-03T00:00:00Z"}, "", true)
	if intent.TargetScope != "none" || intent.Automatic || !intent.RequiresApproval {
		t.Fatalf("unverified asset must not be executable: %#v", intent)
	}
}

func TestActiveDefenseExecutionIsBoundedAndIdempotent(t *testing.T) {
	alert := WazuhAlert{ID: "a1", RuleID: "r1", Severity: "critical", Timestamp: "2026-09-03T00:00:00Z"}
	a := BuildActiveDefenseExecutionIntent(alert, "router-1", true)
	b := BuildActiveDefenseExecutionIntent(alert, "router-1", true)
	if !a.Automatic || a.TargetScope != "ftn-owned-asset" || !a.SnapshotRequired || !a.VerificationRequired {
		t.Fatalf("unexpected execution intent: %#v", a)
	}
	if a.DurationSeconds <= 0 || a.DurationSeconds > 3600 {
		t.Fatalf("unsafe duration: %d", a.DurationSeconds)
	}
	if a.IdempotencyKey == "" || a.IdempotencyKey != b.IdempotencyKey {
		t.Fatal("execution intent must have deterministic idempotency key")
	}
}
