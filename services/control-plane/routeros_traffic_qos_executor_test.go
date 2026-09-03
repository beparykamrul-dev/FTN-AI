package main

import (
	"strings"
	"testing"
	"time"
)

func TestDiffRouterOSTrafficQoSIdempotent(t *testing.T) {
	rule := RouterOSTrafficQoSRule{ServiceID: "pubg", Class: TrafficGaming, DSCP: 46, Priority: 95, PathID: "pop-dhaka"}
	plan := RouterOSTrafficQoSPlan{DeviceID: "r1", Rules: []RouterOSTrafficQoSRule{rule}, RequiresApproval: true}
	diff := DiffRouterOSTrafficQoS(RouterOSTrafficQoSState{Rules: []RouterOSTrafficQoSRule{rule}}, plan)
	if !diff.Empty() { t.Fatalf("expected empty diff: %+v", diff) }
}

func TestDiffRouterOSTrafficQoSDetectsChange(t *testing.T) {
	current := RouterOSTrafficQoSRule{ServiceID: "pubg", Class: TrafficGaming, DSCP: 0, Priority: 40, PathID: "pop-dhaka"}
	desired := RouterOSTrafficQoSRule{ServiceID: "pubg", Class: TrafficGaming, DSCP: 46, Priority: 95, PathID: "pop-dhaka"}
	diff := DiffRouterOSTrafficQoS(RouterOSTrafficQoSState{Rules: []RouterOSTrafficQoSRule{current}}, RouterOSTrafficQoSPlan{DeviceID: "r1", Rules: []RouterOSTrafficQoSRule{desired}})
	if len(diff.Add) != 1 || diff.Add[0].DSCP != 46 { t.Fatalf("unexpected add diff: %+v", diff) }
	if len(diff.Remove) != 0 { t.Fatalf("unexpected remove diff: %+v", diff) }
}

func TestValidateRouterOSTrafficQoSApplyRequiresApprovalSnapshotVerificationRollback(t *testing.T) {
	device := NetworkDevice{ID: "r1", Address: "10.0.0.1", Kind: "core-router", Protocol: "routeros-api", Healthy: true}
	rule := RouterOSTrafficQoSRule{ServiceID: "pubg", Class: TrafficGaming, DSCP: 46, Priority: 95, PathID: "pop-dhaka"}
	diff := RouterOSTrafficQoSDiff{Add: []RouterOSTrafficQoSRule{rule}}
	intent := NetworkExecutionIntent{Device: device, Action: NetworkConfiguration, Approved: false, PrechangeSnapshot: true, VerificationRequired: true, RollbackSafe: true, Timeout: 30 * time.Second}
	if err := ValidateRouterOSTrafficQoSApply(device, intent, diff); err == nil || !strings.Contains(err.Error(), "approval") { t.Fatalf("expected approval error, got %v", err) }
	intent.Approved = true
	if err := ValidateRouterOSTrafficQoSApply(device, intent, diff); err != nil { t.Fatalf("unexpected validation error: %v", err) }
}

func TestValidateRouterOSTrafficQoSApplyRejectsNoChange(t *testing.T) {
	device := NetworkDevice{ID: "r1", Address: "10.0.0.1", Kind: "core-router", Protocol: "routeros-api", Healthy: true}
	intent := NetworkExecutionIntent{Device: device, Action: NetworkConfiguration, Approved: true, PrechangeSnapshot: true, VerificationRequired: true, RollbackSafe: true, Timeout: 30 * time.Second}
	if err := ValidateRouterOSTrafficQoSApply(device, intent, RouterOSTrafficQoSDiff{}); err == nil || !strings.Contains(err.Error(), "no change") { t.Fatalf("expected no-change error, got %v", err) }
}
