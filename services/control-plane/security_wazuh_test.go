package main

import "testing"

func TestWazuhApprovalBoundary(t *testing.T) {
	for _, action := range []WazuhActionClass{WazuhConfiguration, WazuhFirewall, WazuhRouting, WazuhCredential, WazuhDisable, WazuhDelete, WazuhDestructive} {
		if !WazuhActionRequiresApproval(action) {
			t.Fatalf("action %q must require approval", action)
		}
	}
	for _, action := range []WazuhActionClass{WazuhObserve, WazuhCorrelate, WazuhAlertAction, WazuhRecommend, WazuhBackup, WazuhHealthRecover} {
		if WazuhActionRequiresApproval(action) {
			t.Fatalf("action %q should remain automatic", action)
		}
	}
}

func TestNormalizeWazuhSeverity(t *testing.T) {
	if got := NormalizeWazuhSeverity(" HIGH "); got != "high" {
		t.Fatalf("got %q", got)
	}
	if got := NormalizeWazuhSeverity("unknown"); got != "info" {
		t.Fatalf("got %q", got)
	}
}

func TestWazuhAlertHashDeterministic(t *testing.T) {
	a := WazuhAlert{ID: "a1", RuleID: "1001", Severity: "HIGH", AgentRef: "router-1", Timestamp: "2026-09-02T00:00:00Z", Summary: "test"}
	if WazuhAlertHash(a) != WazuhAlertHash(a) {
		t.Fatal("hash must be deterministic")
	}
}
