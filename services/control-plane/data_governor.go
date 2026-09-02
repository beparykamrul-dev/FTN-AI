package main

import (
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

type DataPolicyRequest struct {
	PolicyKey string `json:"policy_key"`
	PolicyJSON json.RawMessage `json:"policy_json"`
}

type DataRequest struct {
	AssetID string `json:"asset_id,omitempty"`
	RequestType string `json:"request_type"`
	RequestJSON json.RawMessage `json:"request_json,omitempty"`
}

func (a *App) dataGovernor(w http.ResponseWriter, r *http.Request) {
	if !requirePermission(a, "data.read", w, r) { return }
	if r.Method != http.MethodGet { method(w, r, http.MethodGet); return }
	if a.db == nil { jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"error":"database_unavailable"}); return }
	rows, err := a.db.Query(r.Context(), `SELECT id, asset_key, name, domain, classification, COALESCE(owner_principal_id::text,''), COALESCE(steward_principal_id::text,''), source_ref, retention_days, status FROM data_assets ORDER BY name`)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_query_failed"}); return }
	defer rows.Close()
	assets := make([]DataAsset, 0)
	for rows.Next() {
		var x DataAsset
		if err := rows.Scan(&x.ID,&x.AssetKey,&x.Name,&x.Domain,&x.Classification,&x.OwnerPrincipalID,&x.StewardPrincipalID,&x.SourceRef,&x.RetentionDays,&x.Status); err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_decode_failed"}); return }
		assets = append(assets, x)
	}
	if err := rows.Err(); err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"data_assets_query_failed"}); return }
	jsonResponse(w, http.StatusOK, map[string]any{"assets":assets})
}

func (a *App) dataRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { method(w, r, http.MethodPost); return }
	var req DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { jsonResponse(w,http.StatusBadRequest,map[string]string{"error":"invalid_json"}); return }
	req.AssetID = strings.TrimSpace(req.AssetID)
	req.RequestType = strings.TrimSpace(req.RequestType)
	if req.RequestType != "export" && req.RequestType != "deletion" && req.RequestType != "retention_override" && req.RequestType != "classification_change" && req.RequestType != "access" { jsonResponse(w,http.StatusBadRequest,map[string]string{"error":"invalid_request_type"}); return }
	permission := "data.read"
	if req.RequestType == "export" { permission = "data.export" }
	if req.RequestType == "deletion" { permission = "data.delete" }
	if req.RequestType == "retention_override" || req.RequestType == "classification_change" || req.RequestType == "access" { permission = "data.change" }
	if !requirePermission(a, permission, w, r) { return }
	if a.db == nil { jsonResponse(w,http.StatusServiceUnavailable,map[string]string{"error":"database_unavailable"}); return }
	var id string
	err := a.db.QueryRow(r.Context(), `INSERT INTO data_requests(asset_id, request_type, request_json) VALUES (NULLIF($1,'')::uuid,$2,COALESCE($3,'{}'::jsonb)) RETURNING id`, req.AssetID, req.RequestType, req.RequestJSON).Scan(&id)
	if err != nil { jsonResponse(w,http.StatusInternalServerError,map[string]string{"error":"data_request_create_failed"}); return }
	a.audit(r, "data.request", id, req.RequestType, req)
	jsonResponse(w,http.StatusAccepted,map[string]any{"id":id,"status":"pending","approval_required":true})
}
