package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type WazuhPipelineRequest struct {
	Alert  WazuhAlert       `json:"alert"`
	Action WazuhActionClass `json:"action,omitempty"`
}

type WazuhCorrelation struct {
	AlertHash         string          `json:"alert_hash"`
	Severity          string          `json:"severity"`
	CorrelationKey    string          `json:"correlation_key"`
	RCA               string          `json:"rca"`
	RecommendedAction WazuhActionClass `json:"recommended_action"`
	RequiresApproval  bool            `json:"requires_approval"`
}

func CorrelateWazuhAlert(a WazuhAlert) WazuhCorrelation {
	severity := NormalizeWazuhSeverity(a.Severity)
	hash := WazuhAlertHash(a)
	recommended := WazuhCorrelate
	rca := "insufficient_evidence"
	switch severity {
	case "critical":
		recommended, rca = WazuhRecommend, "critical_security_event_requires_operator_review"
	case "high":
		recommended, rca = WazuhRecommend, "high_security_event_requires_root_cause_review"
	case "medium":
		recommended, rca = WazuhAlertAction, "medium_security_event_correlated"
	case "low", "info":
		recommended, rca = WazuhCorrelate, "low_confidence_security_event"
	}
	return WazuhCorrelation{AlertHash: hash, Severity: severity, CorrelationKey: "wazuh:" + hash, RCA: rca, RecommendedAction: recommended, RequiresApproval: WazuhActionRequiresApproval(recommended)}
}

func (a *App) createWazuhApproval(r *http.Request, req approvalRequest) (string, error) {
	rc := requestInfo(r)
	if req.ExpiresIn <= 0 { req.ExpiresIn = 3600 }
	if req.ExpiresIn > 86400 { req.ExpiresIn = 86400 }
	hash := requestScopedHash(r, req)
	var id string
	err := a.db.QueryRow(r.Context(), `insert into change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at) values($1,$2,$3,$4,$5,'pending',now()+make_interval(secs=>$6)) on conflict(request_hash) do update set updated_at=now() where change_approvals.tenant_id=$1 returning id::text`, rc.TenantID, rc.PrincipalID, req.Action, req.Resource, hash, req.ExpiresIn).Scan(&id)
	if err != nil { return "", err }
	a.audit(r, "approval.create", id, "pending", map[string]any{"action": req.Action, "resource": req.Resource, "source": "wazuh"})
	return id, nil
}

func (a *App) wazuhPipeline(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "security.alert.ingest", w, r) { return }
	if a.db == nil { jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"error": "database_required"}); return }
	var req WazuhPipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"}); return }
	if strings.TrimSpace(req.Alert.ID) == "" || strings.TrimSpace(req.Alert.AgentRef) == "" || strings.TrimSpace(req.Alert.RuleID) == "" { jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "alert_id_agent_ref_rule_id_required"}); return }
	corr := CorrelateWazuhAlert(req.Alert)
	payload, _ := json.Marshal(map[string]any{"source": "wazuh", "alert_hash": corr.AlertHash, "rule_id": req.Alert.RuleID, "agent_ref": req.Alert.AgentRef, "severity": corr.Severity, "timestamp": req.Alert.Timestamp, "correlation_key": corr.CorrelationKey, "rca": corr.RCA, "recommended_action": corr.RecommendedAction})
	rc := requestInfo(r)
	tx, err := a.db.Begin(r.Context()); if err != nil { jsonResponse(w, 500, map[string]string{"error": "pipeline_begin_failed"}); return }
	defer tx.Rollback(r.Context())
	event, err := appendEventTx(tx, r.Context(), rc.TenantID, "security.wazuh.alert.correlated", rc.CorrelationID, "", corr.CorrelationKey, payload)
	if err != nil { jsonResponse(w, 500, map[string]string{"error": "event_append_failed"}); return }
	if err = tx.Commit(r.Context()); err != nil { jsonResponse(w, 500, map[string]string{"error": "pipeline_commit_failed"}); return }
	result := map[string]any{"event": event, "correlation": corr}
	if req.Action != "" && WazuhActionRequiresApproval(req.Action) {
		body, _ := json.Marshal(map[string]any{"source": "wazuh", "alert_hash": corr.AlertHash, "action": req.Action})
		id, err := a.createWazuhApproval(r, approvalRequest{Action: string(req.Action), Resource: corr.CorrelationKey, Payload: body, ExpiresIn: int((time.Hour).Seconds())})
		if err != nil { jsonResponse(w, 500, map[string]string{"error": "approval_persist_failed"}); return }
		result["approval_required"] = true
		result["approval_id"] = id
	}
	a.audit(r, "security.wazuh.pipeline", corr.CorrelationKey, "correlated", map[string]any{"alert_hash": corr.AlertHash, "severity": corr.Severity, "recommended_action": corr.RecommendedAction})
	jsonResponse(w, http.StatusAccepted, result)
}
