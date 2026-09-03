package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type AIAnalyzeRequest struct {
	Scope   string             `json:"scope,omitempty"`
	Signals map[string]float64 `json:"signals,omitempty"`
	Context map[string]string  `json:"context,omitempty"`
}

type AIFinding struct {
	Code       string  `json:"code"`
	Severity   string  `json:"severity"`
	Confidence float64 `json:"confidence"`
	Summary    string  `json:"summary"`
	Evidence   string  `json:"evidence"`
}

type AIRecommendation struct {
	Code              string `json:"code"`
	Priority          string `json:"priority"`
	Action            string `json:"action"`
	ApprovalRequired  bool   `json:"approval_required"`
	ExecutionAllowed  bool   `json:"execution_allowed"`
}

type AIAnalyzeResponse struct {
	Engine             string             `json:"engine"`
	Mode               string             `json:"mode"`
	Scope              string             `json:"scope"`
	Risk               string             `json:"risk"`
	Findings           []AIFinding        `json:"findings"`
	Recommendations    []AIRecommendation `json:"recommendations"`
	ExecutionPermitted bool               `json:"execution_permitted"`
}

func (a *App) aiStatus(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) { return }
	rc := requestInfo(r)
	jsonResponse(w, http.StatusOK, map[string]any{
		"engine": "ftn-ai-core",
		"version": "1",
		"mode": "decision-support",
		"status": "ready",
		"execution_permitted": false,
		"authenticated_principal": rc.PrincipalID != "",
	})
}

func (a *App) aiAnalyze(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) { return }
	var in AIAnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
		return
	}
	in.Scope = strings.TrimSpace(in.Scope)
	if in.Scope == "" { in.Scope = "network" }
	if len(in.Signals) > 64 || len(in.Context) > 32 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "input_limit_exceeded"})
		return
	}

	resp := analyzeFTNSignals(in)
	a.audit(r, "ai.analyze", in.Scope, "completed", map[string]any{"risk": resp.Risk, "finding_count": len(resp.Findings)})
	jsonResponse(w, http.StatusOK, resp)
}

func analyzeFTNSignals(in AIAnalyzeRequest) AIAnalyzeResponse {
	findings := make([]AIFinding, 0, 8)
	recs := make([]AIRecommendation, 0, 8)
	s := in.Signals

	add := func(f AIFinding, rec AIRecommendation) {
		findings = append(findings, f)
		recs = append(recs, rec)
	}
	if v, ok := s["packet_loss_pct"]; ok && v >= 5 {
		severity, priority := "warning", "high"
		if v >= 15 { severity, priority = "critical", "urgent" }
		add(AIFinding{"packet_loss", severity, 0.96, "Packet loss is elevated", "packet_loss_pct=" + formatAI(v)}, AIRecommendation{"inspect-path", priority, "Inspect interface errors, congestion, and upstream path health", true, false})
	}
	if v, ok := s["latency_ms"]; ok && v >= 100 {
		add(AIFinding{"high_latency", "warning", 0.93, "Network latency is elevated", "latency_ms=" + formatAI(v)}, AIRecommendation{"trace-path", "high", "Trace the affected path and compare alternate healthy routes", true, false})
	}
	if v, ok := s["cpu_pct"]; ok && v >= 85 {
		add(AIFinding{"device_cpu", "warning", 0.94, "Device CPU utilization is high", "cpu_pct=" + formatAI(v)}, AIRecommendation{"inspect-load", "high", "Inspect control-plane load and defer nonessential changes", false, false})
	}
	if v, ok := s["memory_pct"]; ok && v >= 90 {
		add(AIFinding{"device_memory", "critical", 0.95, "Device memory utilization is critically high", "memory_pct=" + formatAI(v)}, AIRecommendation{"protect-device", "urgent", "Protect device stability and investigate memory consumers", true, false})
	}
	if v, ok := s["bgp_established_pct"]; ok && v < 100 {
		priority := "high"
		if v < 75 { priority = "urgent" }
		add(AIFinding{"bgp_degradation", "warning", 0.97, "Not all expected BGP sessions are established", "bgp_established_pct=" + formatAI(v)}, AIRecommendation{"inspect-bgp", priority, "Inspect BGP session state, reachability, and policy changes", true, false})
	}
	if v, ok := s["dns_error_rate_pct"]; ok && v >= 5 {
		add(AIFinding{"dns_errors", "warning", 0.92, "DNS error rate is elevated", "dns_error_rate_pct=" + formatAI(v)}, AIRecommendation{"inspect-dns", "high", "Inspect resolver health, upstream reachability, and recent DNS policy changes", true, false})
	}
	if len(findings) == 0 {
		recs = append(recs, AIRecommendation{"observe", "normal", "Continue telemetry collection; no rule-based anomaly detected", false, false})
	}

	sort.SliceStable(findings, func(i, j int) bool { return findings[i].Confidence > findings[j].Confidence })
	sort.SliceStable(recs, func(i, j int) bool { return recPriority(recs[i].Priority) > recPriority(recs[j].Priority) })
	risk := "low"
	for _, f := range findings { if f.Severity == "critical" { risk = "critical"; break }; if f.Severity == "warning" { risk = "elevated" } }
	return AIAnalyzeResponse{"ftn-ai-core", "1", in.Scope, risk, findings, recs, false}
}

func recPriority(v string) int { switch v { case "urgent": return 3; case "high": return 2; case "normal": return 1 }; return 0 }
func formatAI(v float64) string { return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v), "0"), ".") }
