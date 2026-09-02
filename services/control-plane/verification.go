package main

import (
    "encoding/json"
    "net/http"
    "strings"
)

type verificationRequest struct {
    Success bool            `json:"success"`
    Result  json.RawMessage `json:"result,omitempty"`
    Reason  string          `json:"reason,omitempty"`
}

// verifyJob records post-execution verification. It does not directly change
// network/device state; resource-specific workers remain the execution boundary.
func (a *App) verifyJob(w http.ResponseWriter, r *http.Request) {
    if !method(w,r,http.MethodPost) || !requirePermission(a,"job.verify",w,r) { return }
    if a.db==nil { jsonResponse(w,503,map[string]string{"error":"database_required"}); return }
    id:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/api/v1/jobs/"),"/verify")
    var in verificationRequest
    if json.NewDecoder(r.Body).Decode(&in)!=nil { jsonResponse(w,400,map[string]string{"error":"invalid_verification"}); return }
    if len(in.Result)==0 { in.Result=json.RawMessage(`{}`) }
    rc:=requestInfo(r)
    var status string
    if err:=a.db.QueryRow(r.Context(),`update durable_jobs set verification_payload=$1::jsonb,updated_at=now() where id=$2::uuid and tenant_id=$3::uuid returning status`,string(in.Result),id,rc.TenantID).Scan(&status);err!=nil{jsonResponse(w,404,map[string]string{"error":"job_not_found"});return}
    outcome:="verified";if !in.Success { outcome="verification_failed" }
    a.audit(r,"job.verify",id,outcome,map[string]any{"success":in.Success,"reason":in.Reason,"result":json.RawMessage(in.Result)})
    jsonResponse(w,200,map[string]any{"job_id":id,"status":status,"verified":in.Success})
}
