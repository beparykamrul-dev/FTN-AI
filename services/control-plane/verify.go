package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type verifyRequest struct { WorkerID string `json:"worker_id"`; AttemptNo int `json:"attempt_no"`; Success bool `json:"success"`; Result json.RawMessage `json:"result,omitempty"`; Error string `json:"error,omitempty"` }

// Verification is the final state gate. A privileged approval becomes executed
// only after an authorized verifier records a successful result.
func (a *App) verifyJob(w http.ResponseWriter, r *http.Request) {
	if !method(w,r,http.MethodPost)||!requirePermission(a,"job.verify",w,r){return};if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_required"});return}
	id:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/api/v1/jobs/"),"/verify");var in verifyRequest
	if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil||strings.TrimSpace(in.WorkerID)==""||in.AttemptNo<=0{jsonResponse(w,400,map[string]string{"error":"invalid_verification"});return};if len(in.Result)==0{in.Result=json.RawMessage(`{}`)};if !json.Valid(in.Result){jsonResponse(w,400,map[string]string{"error":"invalid_result"});return}
	rc:=requestInfo(r);if strings.TrimSpace(rc.PrincipalID)==""||strings.TrimSpace(rc.TenantID)==""{jsonResponse(w,401,map[string]string{"error":"principal_required"});return}
	tx,err:=a.db.Begin(r.Context());if err!=nil{jsonResponse(w,500,map[string]string{"error":"verify_begin_failed"});return};defer tx.Rollback(r.Context())
	var status,attemptWorker,attemptStatus string;var approvalID,requestedBy *string
	err=tx.QueryRow(r.Context(),`select j.status,ea.worker_id,ea.status,j.approval_id::text,a.requested_by::text from durable_jobs j join execution_attempts ea on ea.job_id=j.id and ea.attempt_no=$3 left join change_approvals a on a.id=j.approval_id and a.tenant_id=j.tenant_id where j.id=$1::uuid and j.tenant_id=$2::uuid for update`,id,rc.TenantID,in.AttemptNo).Scan(&status,&attemptWorker,&attemptStatus,&approvalID,&requestedBy)
	if err!=nil{if err==pgx.ErrNoRows{jsonResponse(w,404,map[string]string{"error":"job_or_attempt_not_found"})}else{jsonResponse(w,500,map[string]string{"error":"verification_query_failed"})};return}
	if status!="succeeded"||attemptStatus!="succeeded"||attemptWorker!=in.WorkerID{jsonResponse(w,409,map[string]string{"error":"verification_state_conflict"});return};if requestedBy!=nil&&*requestedBy==rc.PrincipalID{jsonResponse(w,409,map[string]string{"error":"verification_separation_required"});return}
	if approvalID!=nil{var approvalStatus string;err=tx.QueryRow(r.Context(),`select status from change_approvals where id=$1::uuid and tenant_id=$2::uuid for update`,*approvalID,rc.TenantID).Scan(&approvalStatus);if err!=nil{jsonResponse(w,500,map[string]string{"error":"approval_state_query_failed"});return};if approvalStatus!="approved"{jsonResponse(w,409,map[string]string{"error":"approval_not_approved"});return}}
	if !in.Success{if _,err=tx.Exec(r.Context(),`update durable_jobs set verification_payload=$2::jsonb,last_error=$3,updated_at=now() where id=$1::uuid and tenant_id=$4::uuid`,id,string(in.Result),in.Error,rc.TenantID);err!=nil{jsonResponse(w,500,map[string]string{"error":"verification_update_failed"});return};if err=tx.Commit(r.Context());err!=nil{jsonResponse(w,500,map[string]string{"error":"verification_commit_failed"});return};a.audit(r,"job.verify",id,"verification_failed",map[string]any{"worker_id":in.WorkerID,"verifier_id":rc.PrincipalID,"attempt":in.AttemptNo,"approval_id":approvalID,"result":json.RawMessage(in.Result),"error":in.Error});jsonResponse(w,409,map[string]any{"job_id":id,"status":"succeeded","verification":"verification_failed"});return}
	if _,err=tx.Exec(r.Context(),`update durable_jobs set verification_payload=$2::jsonb,verified_by=$3::uuid,last_error='',updated_at=now() where id=$1::uuid and tenant_id=$4::uuid and verified_by is null`,id,string(in.Result),rc.PrincipalID,rc.TenantID);err!=nil{jsonResponse(w,500,map[string]string{"error":"verification_update_failed"});return}
	if approvalID!=nil{res,execErr:=tx.Exec(r.Context(),`update change_approvals set status='executed',executed_at=now(),updated_at=now() where id=$1::uuid and tenant_id=$2::uuid and status='approved'`,*approvalID,rc.TenantID);if execErr!=nil{jsonResponse(w,500,map[string]string{"error":"approval_execution_update_failed"});return};if res.RowsAffected()!=1{jsonResponse(w,409,map[string]string{"error":"approval_execution_race"});return}}
	if err=tx.Commit(r.Context());err!=nil{jsonResponse(w,500,map[string]string{"error":"verification_commit_failed"});return};a.audit(r,"job.verify",id,"verified",map[string]any{"worker_id":in.WorkerID,"verifier_id":rc.PrincipalID,"attempt":in.AttemptNo,"approval_id":approvalID,"result":json.RawMessage(in.Result)});jsonResponse(w,200,map[string]any{"job_id":id,"status":"succeeded","verification":"verified","verifier_id":rc.PrincipalID})
}
