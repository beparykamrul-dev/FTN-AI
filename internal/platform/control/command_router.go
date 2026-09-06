package control

import("errors";"strings")
type CommandRouter struct{approvals *ApprovalStore}
func NewCommandRouter(approvals *ApprovalStore)*CommandRouter{return &CommandRouter{approvals:approvals}}
type DispatchResult struct{Status string `json:"status"`;NeedsApproval bool `json:"needs_approval"`;ApprovalID string `json:"approval_id,omitempty"`}
func(r *CommandRouter)Dispatch(serverID string,op Operation,reason string)(DispatchResult,error){if r==nil||r.approvals==nil{return DispatchResult{},errors.New("approval store is required")};serverID=strings.TrimSpace(serverID);reason=strings.TrimSpace(reason);if serverID==""||len(serverID)>256{return DispatchResult{},errors.New("server id is invalid")};if len(reason)>2048{return DispatchResult{},errors.New("reason is too long")};if err:=Validate(op);err!=nil{return DispatchResult{},err};if RequiresApproval(op){id:=serverID+":"+string(op);if len(id)>512{return DispatchResult{},errors.New("approval id is too long")};if err:=r.approvals.Create(Request{ID:id,ServerID:serverID,Operation:op,Reason:reason});err!=nil{return DispatchResult{},err};return DispatchResult{Status:"pending_approval",NeedsApproval:true,ApprovalID:id},nil};return DispatchResult{Status:"ready"},nil}
