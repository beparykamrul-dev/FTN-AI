package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type approvalRequest struct {
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Payload   any    `json:"payload"`
	ExpiresIn int    `json:"expires_in_seconds,omitempty"`
}

type approvalDecision struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

func (a *App) createApproval(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) || !requirePermission(a, "approval.create", w, r) { return }
	if a.db == nil { jsonResponse(w, http.StatusServiceUnavailable, map[string]string{"error":"database_unavailable"}); return }
	var req approvalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Action)=="" || strings.TrimSpace(req.Resource)=="" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error":"action_and_resource_required"}); return
	}
	if req.ExpiresIn <= 0 { req.ExpiresIn = 3600 }
	if req.ExpiresIn > 86400 { req.ExpiresIn = 86400 }
	hash := requestBodyHash(req)
	rc := requestInfo(r)
	var id, status string
	var expires time.Time
	err := a.db.QueryRow(r.Context(), `insert into change_approvals(tenant_id,requested_by,action,resource,request_hash,status,expires_at) values($1,$2,$3,$4,$5,'pending',now()+make_interval(secs=>$6)) on conflict(request_hash) do update set updated_at=now() returning id::text,status,expires_at`, rc.TenantID, rc.PrincipalID, req.Action, req.Resource, hash, req.ExpiresIn).Scan(&id,&status,&expires)
	if err != nil { jsonResponse(w, http.StatusInternalServerError, map[string]string{"error":"approval_persist_failed"}); return }
	a.audit(r,"approval.create",id,"accepted",map[string]any{"action":req.Action,"resource":req.Resource})
	jsonResponse(w,http.StatusAccepted,map[string]any{"approval_id":id,"status":status,"expires_at":expires,"request_hash":hash})
}

func (a *App) listApprovals(w http.ResponseWriter, r *http.Request) {
	if !method(w,r,http.MethodGet) || !requirePermission(a,"approval.read",w,r) { return }
	if a.db==nil { jsonResponse(w,http.StatusServiceUnavailable,map[string]string{"error":"database_unavailable"}); return }
	rc:=requestInfo(r); rows,err:=a.db.Query(r.Context(),`select id::text,action,resource,status,approved_by::text,expires_at,executed_at,created_at from change_approvals where tenant_id=$1 order by created_at desc limit 200`,rc.TenantID)
	if err!=nil {jsonResponse(w,500,map[string]string{"error":"approval_query_failed"});return};defer rows.Close()
	type item struct {ID,Action,Resource,Status string; ApprovedBy *string `json:"approved_by,omitempty"`; ExpiresAt *time.Time `json:"expires_at,omitempty"`; ExecutedAt *time.Time `json:"executed_at,omitempty"`; CreatedAt time.Time `json:"created_at"`}
	out:=[]item{}
	for rows.Next(){var v item;if err:=rows.Scan(&v.ID,&v.Action,&v.Resource,&v.Status,&v.ApprovedBy,&v.ExpiresAt,&v.ExecutedAt,&v.CreatedAt);err!=nil{jsonResponse(w,500,map[string]string{"error":"approval_scan_failed"});return};out=append(out,v)}
	jsonResponse(w,200,map[string]any{"approvals":out})
}

func (a *App) decideApproval(w http.ResponseWriter,r *http.Request){
	if !method(w,r,http.MethodPost) || !requirePermission(a,"approval.decide",w,r){return}
	if a.db==nil{jsonResponse(w,503,map[string]string{"error":"database_unavailable"});return}
	id:=strings.TrimPrefix(r.URL.Path,"/api/v1/approvals/");id=strings.TrimSuffix(id,"/decide");if id==""{jsonResponse(w,400,map[string]string{"error":"approval_id_required"});return}
	var d approvalDecision;if json.NewDecoder(r.Body).Decode(&d)!=nil{jsonResponse(w,400,map[string]string{"error":"invalid_json"});return};d.Decision=strings.ToLower(strings.TrimSpace(d.Decision));if d.Decision!="approve"&&d.Decision!="reject"{jsonResponse(w,400,map[string]string{"error":"invalid_decision"});return}
	rc:=requestInfo(r);tx,err:=a.db.Begin(r.Context());if err!=nil{jsonResponse(w,500,map[string]string{"error":"transaction_failed"});return};defer tx.Rollback(r.Context())
	var current string;err=tx.QueryRow(r.Context(),`select status from change_approvals where id=$1 and tenant_id=$2 for update`,id,rc.TenantID).Scan(&current);if err!=nil{if errors.Is(err,pgx.ErrNoRows){jsonResponse(w,404,map[string]string{"error":"approval_not_found"})}else{jsonResponse(w,500,map[string]string{"error":"approval_query_failed"})};return}
	if current!="pending"{jsonResponse(w,409,map[string]string{"error":"approval_not_pending"});return}
	newStatus:="rejected";if d.Decision=="approve"{newStatus="approved"}
	_,err=tx.Exec(r.Context(),`update change_approvals set status=$1,approved_by=$2,updated_at=now() where id=$3`,newStatus,rc.PrincipalID,id);if err!=nil{jsonResponse(w,500,map[string]string{"error":"approval_update_failed"});return}
	if err=tx.Commit(r.Context());err!=nil{jsonResponse(w,500,map[string]string{"error":"approval_commit_failed"});return}
	a.audit(r,"approval."+d.Decision,id,newStatus,map[string]string{"reason":d.Reason})
	jsonResponse(w,200,map[string]any{"approval_id":id,"status":newStatus})
}
