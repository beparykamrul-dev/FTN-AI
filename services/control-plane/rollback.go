package main

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"

    "github.com/jackc/pgx/v5"
)

type rollbackRequest struct { ApprovalID string `json:"approval_id"`; Reason string `json:"reason,omitempty"` }

// rollbackJob schedules a rollback operation; it never mutates network/device state itself.
// The rollback approval is a distinct approval bound to the exact successful job,
// original approval hash, rollback payload, action, and resource.
func (a *App) rollbackJob(w http.ResponseWriter,r *http.Request){
    if !method(w,r,http.MethodPost)||!requirePermission(a,"job.rollback",w,r){return}
    if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return}
    id:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/api/v1/jobs/"),"/rollback")
    var q rollbackRequest
    if json.NewDecoder(r.Body).Decode(&q)!=nil||strings.TrimSpace(q.ApprovalID)==""{jsonResponse(w,400,map[string]string{"error":"rollback_approval_required"});return}
    rc:=requestInfo(r);tx,err:=a.db.Begin(r.Context());if err!=nil{jsonResponse(w,500,map[string]string{"error":"rollback_begin_failed"});return};defer tx.Rollback(r.Context())

    var jobStatus,action,jobApprovalHash string
    var rollbackPayload []byte
    var approvalStatus,approvalAction,approvalResource,approvalRequestHash string
    err=tx.QueryRow(r.Context(),`select j.status,j.execution_action,coalesce(j.approval_request_hash,''),j.rollback_payload::text,a.status,a.action,a.resource,a.request_hash from durable_jobs j join change_approvals a on a.id=$3::uuid where j.id=$1::uuid and j.tenant_id=$2::uuid`,id,rc.TenantID,q.ApprovalID).Scan(&jobStatus,&action,&jobApprovalHash,&rollbackPayload,&approvalStatus,&approvalAction,&approvalResource,&approvalRequestHash)
    if errors.Is(err,pgx.ErrNoRows){jsonResponse(w,404,map[string]string{"error":"job_or_approval_not_found"});return};if err!=nil{jsonResponse(w,500,map[string]string{"error":"rollback_query_failed"});return}
    if jobStatus!="succeeded"{jsonResponse(w,409,map[string]string{"error":"job_not_successful"});return}
    if approvalStatus!="approved"{jsonResponse(w,409,map[string]string{"error":"rollback_approval_not_approved"});return}
    if strings.TrimSpace(action)==""{jsonResponse(w,409,map[string]string{"error":"job_execution_action_required"});return}
    if len(rollbackPayload)==0{rollbackPayload=[]byte(`{}`)}

    expectedAction:=action+".rollback"
    expectedResource:="durable_job:"+id
    canonical:=approvalRequest{Action:expectedAction,Resource:expectedResource,Payload:map[string]any{"job_id":id,"original_approval_request_hash":jobApprovalHash,"rollback_payload":json.RawMessage(rollbackPayload)}}
    expectedHash:=requestBodyHash(canonical)
    if approvalAction!=expectedAction||approvalResource!=expectedResource||approvalRequestHash!=expectedHash{jsonResponse(w,409,map[string]string{"error":"rollback_approval_binding_mismatch"});return}

    _,err=tx.Exec(r.Context(),`update durable_jobs set status='queued',attempts=0,max_attempts=3,available_at=now(),locked_at=null,locked_by=null,lease_expires_at=null,finished_at=null,last_error='',execution_action=$3,approval_id=$2::uuid,correlation_id=$4,updated_at=now() where id=$1::uuid and tenant_id=$5::uuid and status='succeeded'`,id,q.ApprovalID,expectedAction,rc.CorrelationID,rc.TenantID)
    if err!=nil{jsonResponse(w,500,map[string]string{"error":"rollback_schedule_failed"});return}
    if err=tx.Commit(r.Context());err!=nil{jsonResponse(w,500,map[string]string{"error":"rollback_commit_failed"});return}
    a.audit(r,"job.rollback.schedule",id,"queued",map[string]any{"approval_id":q.ApprovalID,"reason":q.Reason,"previous_action":action,"approval_request_hash":approvalRequestHash})
    jsonResponse(w,202,map[string]any{"job_id":id,"status":"queued","rollback":true,"approval_id":q.ApprovalID,"execution_action":expectedAction})
}
