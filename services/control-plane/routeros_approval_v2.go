package main

import (
	"encoding/json"
	"errors"
	"strings"
)

type RouterOSApprovalBindingV2 struct {
	ApprovalID string `json:"approval_id"`
	DeviceID string `json:"device_id"`
	Plan RouterOSQoSExecutionPlan `json:"plan"`
	Approved bool `json:"approved"`
	Explicit bool `json:"explicit"`
}

func BindRouterOSQoSApprovalV2(plan RouterOSQoSExecutionPlan, approvalID string, explicit bool) (RouterOSApprovalBindingV2, error) {
	approvalID = strings.TrimSpace(approvalID)
	if strings.TrimSpace(plan.DeviceID) == "" { return RouterOSApprovalBindingV2{}, errors.New("routeros_approval_device_id_required") }
	if approvalID == "" || !explicit { return RouterOSApprovalBindingV2{}, errors.New("routeros_explicit_approval_required") }
	if plan.ApplyAllowed { return RouterOSApprovalBindingV2{}, errors.New("routeros_plan_must_remain_apply_closed") }
	if plan.Diff.Empty() { return RouterOSApprovalBindingV2{}, errors.New("routeros_no_change_requires_approval") }
	return RouterOSApprovalBindingV2{ApprovalID: approvalID, DeviceID: plan.DeviceID, Plan: plan, Approved: true, Explicit: true}, nil
}

func EvaluateRouterOSQoSApplyV2(binding RouterOSApprovalBindingV2, current RouterOSQoSExecutionPlan, device NetworkDevice) (NetworkExecutionIntent, error) {
	if !binding.Approved || !binding.Explicit || strings.TrimSpace(binding.ApprovalID) == "" { return NetworkExecutionIntent{}, errors.New("routeros_approval_not_valid") }
	if binding.DeviceID != current.DeviceID || binding.DeviceID != device.ID { return NetworkExecutionIntent{}, errors.New("routeros_approval_device_mismatch") }
	if current.ApplyAllowed || current.Diff.Empty() { return NetworkExecutionIntent{}, errors.New("routeros_apply_plan_invalid") }
	if binding.Plan.DeviceID != current.DeviceID || !routerOSDiffEqualV2(binding.Plan.Diff, current.Diff) { return NetworkExecutionIntent{}, errors.New("routeros_approval_plan_mismatch") }
	intent := RouterOSTrafficQoSAction(device)
	intent.Action = NetworkConfiguration
	intent.ApprovalID = binding.ApprovalID
	intent.Approved = true
	intent.Explicit = true
	return intent, nil
}

func routerOSDiffEqualV2(a, b RouterOSDiff) bool {
	x, errA := json.Marshal(a); y, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(x) == string(y)
}
