package main

import "testing"

func TestBindRouterOSQoSApprovalRequiresExplicitApproval(t *testing.T) {
	plan := RouterOSQoSExecutionPlan{DeviceID: "r1", Diff: RouterOSDiff{Adds: []RouterOSQoSState{{ServiceID: "whatsapp", PathID: "p1"}}}}
	if _, err := BindRouterOSQoSApproval(plan, "", false); err == nil { t.Fatal("expected approval requirement") }
	binding, err := BindRouterOSQoSApproval(plan, "apr-1", true)
	if err != nil || !binding.Approved || !binding.Explicit { t.Fatalf("unexpected binding: %+v %v", binding, err) }
}

func TestEvaluateRouterOSQoSApplyRejectsChangedPlan(t *testing.T) {
	plan := RouterOSQoSExecutionPlan{DeviceID: "r1", Diff: RouterOSDiff{Adds: []RouterOSQoSState{{ServiceID: "whatsapp", PathID: "p1"}}}}
	binding, err := BindRouterOSQoSApproval(plan, "apr-1", true)
	if err != nil { t.Fatal(err) }
	changed := RouterOSQoSExecutionPlan{DeviceID: "r1", Diff: RouterOSDiff{Adds: []RouterOSQoSState{{ServiceID: "telegram", PathID: "p1"}}}}
	if _, err := EvaluateRouterOSQoSApply(binding, changed, NetworkDevice{ID:"r1", Kind:"router", Address:"10.0.0.1", Protocol:"routeros-api", Healthy:true}); err == nil { t.Fatal("expected plan mismatch") }
}

func TestEvaluateRouterOSQoSApplyProducesApprovedIntent(t *testing.T) {
	plan := RouterOSQoSExecutionPlan{DeviceID: "r1", Diff: RouterOSDiff{Adds: []RouterOSQoSState{{ServiceID: "whatsapp", PathID: "p1"}}}}
	binding, err := BindRouterOSQoSApproval(plan, "apr-1", true)
	if err != nil { t.Fatal(err) }
	intent, err := EvaluateRouterOSQoSApply(binding, plan, NetworkDevice{ID:"r1", Kind:"router", Address:"10.0.0.1", Protocol:"routeros-api", Healthy:true})
	if err != nil { t.Fatal(err) }
	if !intent.Approved || !intent.Explicit || intent.ApprovalID != "apr-1" { t.Fatalf("unexpected intent: %+v", intent) }
	if !intent.PreSnapshot || !intent.PostVerify || !intent.RollbackSafe { t.Fatalf("safety boundary missing: %+v", intent) }
}
