package main

import "testing"

func TestBuildFTNFailoverIntent(t *testing.T) {
	intent, err := BuildFTNFailoverIntent(FTNCoreFailoverDecision{
		ActiveNode: "core-b", Failover: true, Reason: "alternate_core_healthy",
	})
	if err != nil { t.Fatal(err) }
	if intent.TargetNode != "core-b" || intent.DecisionHash == "" {
		t.Fatalf("unexpected intent: %+v", intent)
	}
}

func TestBuildFTNFailoverIntentRejectsNoop(t *testing.T) {
	if _, err := BuildFTNFailoverIntent(FTNCoreFailoverDecision{ActiveNode: "core-a", Failover: false}); err == nil {
		t.Fatal("expected no-op rejection")
	}
}
