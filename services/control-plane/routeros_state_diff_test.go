package main

import (
	"testing"
	"time"
)

func qosRule(service, path string, class TrafficClass, dscp, priority uint8) RouterOSQoSState {
	return RouterOSQoSState{ServiceID: service, PathID: path, Class: class, DSCP: dscp, Priority: priority}
}

func TestDiffRouterOSQoSStateNoOp(t *testing.T) {
	now := time.Now()
	rule := qosRule("whatsapp", "path-b", "voice", 46, 90)
	diff, err := DiffRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1", Rules: []RouterOSQoSState{rule}, CapturedAt: now},
		RouterOSDesiredState{DeviceID: "r1", Rules: []RouterOSQoSState{rule}},
	)
	if err != nil { t.Fatal(err) }
	if !diff.Empty() { t.Fatalf("expected no-op diff, got %#v", diff) }
}

func TestDiffRouterOSQoSStateAddChangeRemove(t *testing.T) {
	before := qosRule("whatsapp", "path-a", "voice", 46, 90)
	change := qosRule("telegram", "path-a", "realtime", 46, 90)
	after := qosRule("whatsapp", "path-a", "voice", 34, 95)
	diff, err := DiffRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1", Rules: []RouterOSQoSState{before, change}},
		RouterOSDesiredState{DeviceID: "r1", Rules: []RouterOSQoSState{after, qosRule("imo", "path-z", "voice", 46, 92)}},
	)
	if err != nil { t.Fatal(err) }
	if len(diff.Adds) != 1 || diff.Adds[0].ServiceID != "imo" { t.Fatalf("unexpected adds: %#v", diff.Adds) }
	if len(diff.Changes) != 1 || diff.Changes[0].Before.ServiceID != "whatsapp" || diff.Changes[0].After.DSCP != 34 { t.Fatalf("unexpected changes: %#v", diff.Changes) }
	if len(diff.Removes) != 1 || diff.Removes[0].ServiceID != "telegram" { t.Fatalf("unexpected removes: %#v", diff.Removes) }
}

func TestNormalizeRouterOSQoSStateDeterministic(t *testing.T) {
	rules, err := NormalizeRouterOSQoSState([]RouterOSQoSState{
		qosRule("telegram", "path-z", "realtime", 46, 90),
		qosRule("whatsapp", "path-a", "voice", 46, 90),
	})
	if err != nil { t.Fatal(err) }
	if len(rules) != 2 || rules[0].ServiceID != "telegram" || rules[1].ServiceID != "whatsapp" { t.Fatalf("unexpected order: %#v", rules) }
}

func TestNormalizeRouterOSQoSStateRejectsInvalidInput(t *testing.T) {
	cases := []RouterOSQoSState{
		{ServiceID: "", PathID: "path-a"},
		{ServiceID: "whatsapp", PathID: ""},
		{ServiceID: "whatsapp", PathID: "path-a", DSCP: 64},
		{ServiceID: "not-a-managed-service", PathID: "path-a"},
	}
	for _, rule := range cases {
		if _, err := NormalizeRouterOSQoSState([]RouterOSQoSState{rule}); err == nil { t.Fatalf("expected rejection for %#v", rule) }
	}
}

func TestDiffRouterOSQoSStateRejectsDeviceMismatch(t *testing.T) {
	rule := qosRule("whatsapp", "path-a", "voice", 46, 90)
	if _, err := DiffRouterOSQoSState(
		RouterOSSnapshot{DeviceID: "r1", Rules: []RouterOSQoSState{rule}},
		RouterOSDesiredState{DeviceID: "r2", Rules: []RouterOSQoSState{rule}},
	); err == nil { t.Fatal("expected device mismatch error") }
}

func TestNormalizeRouterOSQoSStateRejectsDuplicates(t *testing.T) {
	rule := qosRule("whatsapp", "path-a", "voice", 46, 90)
	if _, err := NormalizeRouterOSQoSState([]RouterOSQoSState{rule, rule}); err == nil { t.Fatal("expected duplicate rejection") }
}
