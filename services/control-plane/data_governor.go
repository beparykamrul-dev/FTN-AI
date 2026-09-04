package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type DataAsset struct {
	ID string `json:"id"`
	AssetKey string `json:"asset_key"`
	Name string `json:"name"`
	Domain string `json:"domain"`
	Classification string `json:"classification"`
	OwnerPrincipalID string `json:"owner_principal_id,omitempty"`
	StewardPrincipalID string `json:"steward_principal_id,omitempty"`
	SourceRef string `json:"source_ref,omitempty"`
	RetentionDays *int `json:"retention_days,omitempty"`
	Status string `json:"status"`
}

type DataRequest struct {
	AssetID string `json:"asset_id,omitempty"`
	RequestType string `json:"request_type"`
	RequestJSON json.RawMessage `json:"request_json,omitempty"`
}

var forbiddenDataKeys = map[string]struct{}{
	"password": {}, "passwd": {}, "secret": {}, "token": {}, "access_token": {}, "refresh_token": {},
	"api_key": {}, "apikey": {}, "private_key": {}, "privatekey": {}, "client_secret": {}, "qkd_key": {}, "raw_key": {},
}

func containsSensitiveKey(v any) bool {
	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(k), "-", "_"))
			if _, ok := forbiddenDataKeys[key]; ok { return true }
			if containsSensitiveKey(child) { return true }
		}
	case []any:
		for _, child := range x { if containsSensitiveKey(child) { return true } }
	}
	return false
}

func decodeGovernedJSON(raw json.RawMessage) (map[string]any, error) {
	if len(bytes.TrimSpace(raw)) == 0 { return map[string]any{}, nil }
	var v map[string]any
	if err := json.Unmarshal(raw, &v); err != nil { return nil, err }
	if containsSensitiveKey(v) { return nil, errSensitiveGovernanceMetadata }
	return v, nil
}

func (a *App) dataGovernor(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) { return }
	if !requirePermission(a, "data.read", w, r) { return }
	if a.db == nil { jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"error":"database_unavailable"}); return }
	rc := requestInfo(r)
	rows, err := a.db.Query(r.Context(), `SELECT id, asset_key, name, domain, classification, COALESCE(owner_principal_id::text,''), COALESCE(steward_principal_id::text,''), source_ref, retention_days, status FROM data_assets WHERE tenant_id=$1::uuid ORDER BY name`, rc.TenantID)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_query_failed"}); return }
	defer rows.Close()
	assets := make([]DataAsset, 0)
	for rows.Next() {
		var x DataAsset
		if err := rows.Scan(&x.ID, &x.AssetKey, &x.Name, &x.Domain, &x.Classification, &x.OwnerPrincipalID, &x.StewardPrincipalID, &x.SourceRef, &x.RetentionDays, &x.Status); err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_decode_failed"}); return }
		assets = append(assets, x)
	}
	if err := rows.Err(); err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_query_failed"}); return }
	jsonResponse(w, http.StatusOK, map[string]any{"assets": assets})
}

func (a *App) dataRequest(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) { return }
	var req DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"invalid_json"}); return }
	req.AssetID = strings.TrimSpace(req.AssetID)
	req.RequestType = strings.TrimSpace(req.RequestType)
	if req.RequestType != "export" && req.RequestType != "deletion" && req.RequestType != "retention_override" && req.RequestType != "classification_change" && req.RequestType != "access" { jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"invalid_request_type"}); return }
	permission := "data.read"
	if req.RequestType == "export" { permission = "data.export" }
	if req.RequestType == "deletion" { permission = "data.delete" }
	if req.RequestType == "retention_override" || req.RequestType == "classification_change" || req.RequestType == "access" { permission = "data.change" }
	if !requirePermission(a, permission, w, r) { return }
	payload, err := decodeGovernedJSON(req.RequestJSON)
	if err != nil { jsonResponse(w, http.StatusUnprocessableEntity, map[string]string{"error":"sensitive_metadata_rejected"}); return }
	if a.db == nil { jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"error":"database_unavailable"}); return }
	rc := requestInfo(r)
	if req.AssetID != "" {
		var exists bool
		if err := a.db.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM data_assets WHERE id=$1::uuid AND tenant_id=$2::uuid)`, req.AssetID, rc.TenantID).Scan(&exists); err != nil || !exists { jsonResponse(w, http.StatusNotFound, map[string]string{"error":"asset_not_found"}); return }
	}
	approvalAction := "data." + req.RequestType
	approvalResource := "data-governor/request"
	approvalPayload := map[string]any{"asset_id": req.AssetID, "request_type": req.RequestType, "request_json": payload}
	hash := requestScopedHash(r, approvalPayload)
	var approvalID string
	err = a.db.QueryRow(r.Context(), `INSERT INTO change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at) VALUES($1,$2,$3,$4,$5,'pending',now()+interval '1 hour') ON CONFLICT(tenant_id,request_hash) DO UPDATE SET updated_at=now() RETURNING id::text`, rc.TenantID, rc.PrincipalID, approvalAction, approvalResource, hash).Scan(&approvalID)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"approval_persist_failed"}); return }
	encoded, err := json.Marshal(payload)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"request_encode_failed"}); return }
	var id, status string
	err = a.db.QueryRow(r.Context(), `INSERT INTO data_requests(tenant_id,asset_id,request_type,status,requested_by,approval_id,request_json,request_hash) VALUES($1,NULLIF($2,'')::uuid,$3,'pending',$4,$5::uuid,$6::jsonb,$7) ON CONFLICT(tenant_id,request_hash) DO UPDATE SET updated_at=data_requests.updated_at RETURNING id::text,status`, rc.TenantID, req.AssetID, req.RequestType, rc.PrincipalID, approvalID, string(encoded), hash).Scan(&id, &status)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_request_create_failed"}); return }
	a.audit(r, "data.request", id, status, map[string]any{"approval_id": approvalID, "request_type": req.RequestType})
	jsonResponse(w, http.StatusAccepted, map[string]any{"id": id, "status": status, "approval_required": true, "approval_id": approvalID, "request_hash": hash})
}

var errSensitiveGovernanceMetadata = sensitiveMetadataError{}
type sensitiveMetadataError struct{}
func (s sensitiveMetadataError) Error() string { return "sensitive governance metadata" }
