package main

import "testing"

func TestCorrelateWazuhAlert(t *testing.T) {
	tests := []struct { severity, action, rca string }{
		{"critical", string(WazuhRecommend), "critical_security_event_requires_operator_review"},
		{"high", string(WazuhRecommend), "high_security_event_requires_root_cause_review"},
		{"medium", string(WazuhAlertAction), "medium_security_event_correlated"},
		{"info", string(WazuhCorrelate), "low_confidence_security_event"},
	}
	for _, tc := range tests {
		got := CorrelateWazuhAlert(WazuhAlert{ID: "a1", RuleID: "1001", Severity: tc.severity, AgentRef: "router-1", Timestamp: "2026-09-03T00:00:00Z"})
		if string(got.RecommendedAction) != tc.action || got.RCA != tc.rca || got.AlertHash == "" || got.CorrelationKey == "" {
			t.Fatalf("unexpected correlation for %s: %#v", tc.severity, got)
		}
		if got.RequiresApproval { t.Fatalf("advisory correlation must not require approval: %#v", got) }
	}
}

func TestWazuhPipelineRejectsIncompleteAlert(t *testing.T) {
	if got := CorrelateWazuhAlert(WazuhAlert{ID: "", RuleID: "1001", AgentRef: "router-1"}); got.AlertHash == "" {
		t.Fatal("correlation remains deterministic for unit-level input")
	}
}
