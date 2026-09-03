package main

import (
	"errors"
	"strings"
)

// RouterOSApprovalBinding binds an approval to the exact planned diff. The
// approval cannot be reused for another device, plan, or changed diff.
type RouterOSApprovalBinding struct {
	ApprovalID string `json:"approval_id"`
	DeviceID string `json:"device_id"`
	Plan RouterOSQoSExecutionPlan `json:"plan"`
	Approved bool `json:"approved"`
	Explicit bool `json:"explicit"`
}

func BindRouterOSQoSApproval(plan RouterOSQoSExecutionPlan, approvalID string, explicit bool) (RouterOSApprovalBinding, error) {
	approvalID = strings.TrimSpace(approvalID)
	if strings.TrimSpace(plan.DeviceID) == "" { return RouterOSApprovalBinding{}, errors.New("routeros_approval_device_id_required") }
	if approvalID == "" || !explicit { return RouterOSApprovalBinding{}, errors.New("routeros_explicit_approval_required") }
	if plan.ApplyAllowed { return RouterOSApprovalBinding{}, errors.New("routeros_plan_must_remain_apply_closed") }
	if plan.Diff.Empty() { return RouterOSApprovalBinding{}, errors.New("routeros_no_change_requires_approval") }
	return RouterOSApprovalBinding{ApprovalID: approvalID, DeviceID: plan.DeviceID, Plan: plan, Approved: true, Explicit: true}, nil
}

// EvaluateRouterOSQoSApply is the final approval boundary. It validates the
// binding against the current plan and returns an execution intent only when
// the exact approved plan is supplied. It performs no network I/O.
func EvaluateRouterOSQoSApply(binding RouterOSApprovalBinding, current RouterOSQoSExecutionPlan, device NetworkDevice) (NetworkExecutionIntent, error) {
	if !binding.Approved || !binding.Explicit || strings.TrimSpace(binding.ApprovalID) == "" { return NetworkExecutionIntent{}, errors.New("routeros_approval_not_valid") }
	if binding.DeviceID != current.DeviceID || binding.DeviceID != device.ID { return NetworkExecutionIntent{}, errors.New("routeros_approval_device_mismatch") }
	if current.ApplyAllowed || current.Diff.Empty() { return NetworkExecutionIntent{}, errors.New("routeros_apply_plan_invalid") }
	if binding.Plan.DeviceID != current.DeviceID || !routerOSDiffEqual(binding.Plan.Diff, current.Diff) { return NetworkExecutionIntent{}, errors.New("routeros_approval_plan_mismatch") }
	intent := RouterOSTrafficQoSAction(device)
	intent.Action = NetworkConfiguration
	intent.ApprovalID = binding.ApprovalID
	intent.Approved = true
	intent.Explicit = true
	return intent, nil
}

func routerOSDiffEqual(a, b RouterOSDiff) bool {
	return string(mustJSON(a)) == string(mustJSON(b))
}

func mustJSON(v any) []byte {
	b, _ := jsonMarshal(v)
	return b
}
