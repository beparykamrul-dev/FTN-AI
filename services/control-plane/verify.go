package main

import (
    "encoding/json"
    "net/http"
    "strings"
)

type verifyRequest struct {
    WorkerID  string          `json:"worker_id"`
    AttemptNo int             `json:"attempt_no"`
    Success   bool            `json:"success"`
    Result    json.RawMessage `json:"result,omitempty"`
    Error     string          `json:"error,omitempty"`
}

// verifyJob records an explicit post-execution verification result. It does not
// mutate network/device state; resource-specific workers remain responsible for
// performing checks and rollback actions through the approved execution path.
func (a *App) verifyJob(w http.ResponseWriter, r *http.Request) {
    if !method(w, r, http.MethodPost) || !requirePermission(a, "job.verify", w, r) { return }
    if a.db == nil { jsonResponse(w, 503, map[string]string{"error":"database_required"}); return }
    id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/"), "/verify")
    var in verifyRequest
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil || strings.TrimSpace(in.WorkerID) == "" || in.AttemptNo <= 0 {
        jsonResponse(w, 400, map[string]string{"error":"invalid_verification"}); return
    }
    if len(in.Result) == 0 { in.Result = json.RawMessage(`{}`) }
    rc := requestInfo(r)
    tx, err := a.db.Begin(r.Context()); if err != nil { jsonResponse(w,500,map[string]string{"error":"verify_begin_failed"}); return }
    defer tx.Rollback(r.Context())

    var status, lockedBy string
    var attempt int
    if err = tx.QueryRow(r.Context(), `select status,coalesce(locked_by,''),attempts from durable_jobs where id=$1::uuid and tenant_id=$2::uuid for update`, id, rc.TenantID).Scan(&status,&lockedBy,&attempt); err != nil {
        jsonResponse(w,404,map[string]string{"error":"job_not_found"}); return
    }
    if status != "succeeded" || lockedBy != "" || attempt != in.AttemptNo {
        jsonResponse(w,409,map[string]string{"error":"verification_state_conflict"}); return
    }
    verificationStatus := "verified"
    if !in.Success { verificationStatus = "verification_failed" }
    if _, err = tx.Exec(r.Context(), `update durable_jobs set verification_payload=$2::jsonb,last_error=$3,updated_at=now() where id=$1::uuid and tenant_id=$4::uuid`, id, string(in.Result), in.Error, rc.TenantID); err != nil {
        jsonResponse(w,500,map[string]string{"error":"verification_update_failed"}); return
    }
    if err = tx.Commit(r.Context()); err != nil { jsonResponse(w,500,map[string]string{"error":"verification_commit_failed"}); return }
    a.audit(r,"job.verify",id,verificationStatus,map[string]any{"worker_id":in.WorkerID,"attempt":in.AttemptNo,"result":json.RawMessage(in.Result),"error":in.Error})
    jsonResponse(w,200,map[string]any{"job_id":id,"status":"succeeded","verification":verificationStatus})
}
