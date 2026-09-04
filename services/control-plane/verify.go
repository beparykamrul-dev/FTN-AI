package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type verifyRequest struct {
	WorkerID  string          `json:"worker_id"`
	AttemptNo int             `json:"attempt_no"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
}

// Verification is the final state gate. A privileged approval becomes executed
// only after the claimed worker's successful execution has been independently
// verified. This endpoint records state; it never performs device mutations.
func (a *App) verifyJob(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "job.verify", w, r) {
		return
	}
	if a.db == nil {
		jsonResponse(w, 503, map[string]string{"error": "database_required"})
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/"), "/verify")
	var in verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || strings.TrimSpace(in.WorkerID) == "" || in.AttemptNo <= 0 {
		jsonResponse(w, 400, map[string]string{"error": "invalid_verification"})
		return
	}
	if len(in.Result) == 0 {
		in.Result = json.RawMessage(`{}`)
	}
	if !json.Valid(in.Result) {
		jsonResponse(w, 400, map[string]string{"error": "invalid_result"})
		return
	}
	rc := requestInfo(r)
	tx, err := a.db.Begin(r.Context())
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "verify_begin_failed"})
		return
	}
	defer tx.Rollback(r.Context())

	var status string
	var attemptWorker, attemptStatus string
	var approvalID *string
	err = tx.QueryRow(r.Context(), `select j.status,ea.worker_id,ea.status,j.approval_id::text from durable_jobs j join execution_attempts ea on ea.job_id=j.id and ea.attempt_no=$3 where j.id=$1::uuid and j.tenant_id=$2::uuid for update`, id, rc.TenantID, in.AttemptNo).Scan(&status, &attemptWorker, &attemptStatus, &approvalID)
	if err != nil {
		if err == pgx.ErrNoRows {
			jsonResponse(w, 404, map[string]string{"error": "job_or_attempt_not_found"})
		} else {
			jsonResponse(w, 500, map[string]string{"error": "verification_query_failed"})
		}
		return
	}
	if status != "succeeded" || attemptStatus != "succeeded" || attemptWorker != in.WorkerID {
		jsonResponse(w, 409, map[string]string{"error": "verification_state_conflict"})
		return
	}

	if !in.Success {
		if _, err = tx.Exec(r.Context(), `update durable_jobs set verification_payload=$2::jsonb,last_error=$3,updated_at=now() where id=$1::uuid`, id, string(in.Result), in.Error); err != nil {
			jsonResponse(w, 500, map[string]string{"error": "verification_update_failed"})
			return
		}
		if err = tx.Commit(r.Context()); err != nil {
			jsonResponse(w, 500, map[string]string{"error": "verification_commit_failed"})
			return
		}
		a.audit(r, "job.verify", id, "verification_failed", map[string]any{"worker_id": in.WorkerID, "attempt": in.AttemptNo, "approval_id": approvalID, "result": json.RawMessage(in.Result), "error": in.Error})
		jsonResponse(w, 409, map[string]any{"job_id": id, "status": "succeeded", "verification": "verification_failed"})
		return
	}

	if _, err = tx.Exec(r.Context(), `update durable_jobs set verification_payload=$2::jsonb,last_error='',updated_at=now() where id=$1::uuid`, id, string(in.Result)); err != nil {
		jsonResponse(w, 500, map[string]string{"error": "verification_update_failed"})
		return
	}
	if approvalID != nil {
		if _, err = tx.Exec(r.Context(), `update change_approvals set status='executed',executed_at=now(),updated_at=now() where id=$1::uuid and tenant_id=$2::uuid and status='approved'`, *approvalID, rc.TenantID); err != nil {
			jsonResponse(w, 500, map[string]string{"error": "approval_execution_update_failed"})
			return
		}
	}
	if err = tx.Commit(r.Context()); err != nil {
		jsonResponse(w, 500, map[string]string{"error": "verification_commit_failed"})
		return
	}
	a.audit(r, "job.verify", id, "verified", map[string]any{"worker_id": in.WorkerID, "attempt": in.AttemptNo, "approval_id": approvalID, "result": json.RawMessage(in.Result)})
	jsonResponse(w, 200, map[string]any{"job_id": id, "status": "succeeded", "verification": "verified"})
}
