package main

import (
	"testing"
	"time"
)

func TestVerifyRouterOSQoSStateExactMatch(t *testing.T) {
	rules := []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 46, Priority: 90, PathID: "p1"}}
	got, err := VerifyRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1", Rules: rules, CapturedAt: time.Now().UTC()},
		RouterOSDesiredState{DeviceID: "r1", Rules: rules},
	)
	if err != nil { t.Fatal(err) }
	if !got.Verified || got.Drift || !got.Diff.Empty() { t.Fatalf("unexpected verification: %+v", got) }
}

func TestVerifyRouterOSQoSStateDetectsDrift(t *testing.T) {
	got, err := VerifyRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1", Rules: []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 46, Priority: 90, PathID: "p1"}}},
		RouterOSDesiredState{DeviceID: "r1", Rules: []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 34, Priority: 90, PathID: "p1"}}},
	)
	if err != nil { t.Fatal(err) }
	if got.Verified || !got.Drift || len(got.Diff.Changes) != 1 { t.Fatalf("expected one drift change: %+v", got) }
}

func TestVerifyRouterOSQoSStateFailsClosedOnMismatch(t *testing.T) {
	if _, err := VerifyRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1"},
		RouterOSDesiredState{DeviceID: "r2"},
	); err == nil { t.Fatal("expected device mismatch") }
}

func TestBuildRouterOSQoSRollbackPlanTargetsPreChangeState(t *testing.T) {
	pre := RouterOSSnapshot{DeviceID: "r1", CapturedAt: time.Now().UTC(), Rules: []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 46, Priority: 90, PathID: "p1"}}}
	cur := RouterOSSnapshot{DeviceID: "r1", CapturedAt: time.Now().UTC(), Rules: []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 34, Priority: 90, PathID: "p1"}}}
	plan, err := BuildRouterOSQoSRollbackPlan(pre, cur)
	if err != nil { t.Fatal(err) }
	if plan.ApplyAllowed { t.Fatal("rollback plan must never authorize apply") }
	if !plan.RequiresApproval { t.Fatal("non-empty rollback must require approval") }
	if len(plan.Diff.Changes) != 1 || plan.Target.Rules[0].DSCP != 46 { t.Fatalf("unexpected rollback plan: %+v", plan) }
}

func TestBuildRouterOSQoSRollbackPlanNoOpStillClosed(t *testing.T) {
	s := RouterOSSnapshot{DeviceID: "r1", Rules: []RouterOSQoSState{{ServiceID: "whatsapp", Class: TrafficClassVoice, DSCP: 46, Priority: 90, PathID: "p1"}}}
	plan, err := BuildRouterOSQoSRollbackPlan(s, s)
	if err != nil { t.Fatal(err) }
	if plan.ApplyAllowed { t.Fatal("no-op rollback plan must remain closed") }
	if plan.RequiresApproval { t.Fatal("no-op rollback does not require approval") }
	if !plan.Diff.Empty() { t.Fatal("expected no-op rollback diff") }
}
