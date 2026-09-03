package main

import "testing"

func TestAnalyzeFTNSignalsDetectsCriticalNetworkRisk(t *testing.T) {
	out := analyzeFTNSignals(AIAnalyzeRequest{Scope: "core", Signals: map[string]float64{
		"packet_loss_pct": 18,
		"latency_ms": 140,
		"bgp_established_pct": 60,
	}})
	if out.Risk != "critical" { t.Fatalf("risk=%q", out.Risk) }
	if len(out.Findings) != 3 { t.Fatalf("findings=%d", len(out.Findings)) }
	if out.ExecutionPermitted { t.Fatal("AI analysis must never directly permit execution") }
	for _, rec := range out.Recommendations { if rec.ExecutionAllowed { t.Fatalf("recommendation %q unexpectedly executable", rec.Code) } }
}

func TestAnalyzeFTNSignalsHealthyObservation(t *testing.T) {
	out := analyzeFTNSignals(AIAnalyzeRequest{Signals: map[string]float64{
		"packet_loss_pct": 0.2,
		"latency_ms": 20,
		"cpu_pct": 40,
		"bgp_established_pct": 100,
	}})
	if out.Risk != "low" { t.Fatalf("risk=%q", out.Risk) }
	if len(out.Findings) != 0 { t.Fatalf("findings=%d", len(out.Findings)) }
	if len(out.Recommendations) != 1 || out.Recommendations[0].Code != "observe" { t.Fatalf("recommendations=%+v", out.Recommendations) }
}

func TestFormatAI(t *testing.T) {
	if got := formatAI(12.5); got != "12.5" { t.Fatalf("got=%q", got) }
	if got := formatAI(12); got != "12" { t.Fatalf("got=%q", got) }
}
