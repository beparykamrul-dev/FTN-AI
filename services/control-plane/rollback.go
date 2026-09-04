package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type rollbackRequest struct {
	ApprovalID string `json:"approval_id"`
	Reason     string `json:"reason,omitempty"`
}

// rollbackJob schedules the rollback operation; it never mutates network/device state itself.
// A separate approval is mandatory, must explicitly target this job, and must bind
// the exact rollback payload that will be executed.
func (a *App) rollbackJob(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "job.rollback", w, r) {
		return
	}
	if a.db == nil {
		jsonResponse(w, 503, map[string]string{"error": "database_required"})
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/"), "/rollback")
	var q rollbackRequest
	if json.NewDecoder(r.Body).Decode(&q) != nil || strings.TrimSpace(q.ApprovalID) == "" {
		jsonResponse(w, 400, map[string]string{"error": "rollback_approval_required"})
		return
	}
	rc := requestInfo(r)
	tx, err := a.db.Begin(r.Context())
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "rollback_begin_failed"})
		return
	}
	defer tx.Rollback(r.Context())

	var jobStatus, approvalStatus, action, approvalAction, approvalResource, payloadHash string
	var rollbackPayload []byte
	err = tx.QueryRow(r.Context(), `
		select j.status,a.status,j.execution_action,j.rollback_payload::text,
		       a.action,a.resource,coalesce(a.payload_hash,'')
		from durable_jobs j
		join change_approvals a on a.id=$3::uuid
		where j.id=$1::uuid and j.tenant_id=$2::uuid and a.tenant_id=$2::uuid
		  and a.action='job.rollback' and a.resource=$1::text
		  and a.status='approved' and (a.expires_at is null or a.expires_at>now())
		for update of j,a
	`, id, rc.TenantID, q.ApprovalID).Scan(&jobStatus, &approvalStatus, &action, &rollbackPayload, &approvalAction, &approvalResource, &payloadHash)
	if errors.Is(err, pgx.ErrNoRows) {
		jsonResponse(w, 404, map[string]string{"error": "job_or_approval_not_found"})
		return
	}
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "rollback_query_failed"})
		return
	}
	if jobStatus != "succeeded" {
		jsonResponse(w, 409, map[string]string{"error": "job_not_successful"})
		return
	}
	if approvalStatus != "approved" || approvalAction != "job.rollback" || approvalResource != id {
		jsonResponse(w, 409, map[string]string{"error": "rollback_approval_not_bound"})
		return
	}
	if len(rollbackPayload) == 0 {
		rollbackPayload = []byte(`{}`)
	}
	if payloadHash == "" || payloadHash != requestBodyHash(json.RawMessage(rollbackPayload)) {
		jsonResponse(w, 409, map[string]string{"error": "rollback_payload_mismatch"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		update durable_jobs
		set status='queued',attempts=0,max_attempts=3,available_at=now(),locked_at=null,
		    locked_by=null,finished_at=null,last_error='',payload=$3::jsonb,
		    execution_action=case
		      when execution_action='' then 'rollback' else execution_action||'.rollback' end,
		    approval_id=$2::uuid,rollback_payload=$3::jsonb,correlation_id=$4,updated_at=now()
		where id=$1::uuid and tenant_id=$5::uuid
	`, id, q.ApprovalID, string(rollbackPayload), rc.CorrelationID, rc.TenantID)
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "rollback_schedule_failed"})
		return
	}
	if err = tx.Commit(r.Context()); err != nil {
		jsonResponse(w, 500, map[string]string{"error": "rollback_commit_failed"})
		return
	}
	a.audit(r, "job.rollback.schedule", id, "queued", map[string]any{"approval_id": q.ApprovalID, "reason": q.Reason, "previous_action": action})
	jsonResponse(w, 202, map[string]any{"job_id": id, "status": "queued", "rollback": true, "approval_id": q.ApprovalID})
}
