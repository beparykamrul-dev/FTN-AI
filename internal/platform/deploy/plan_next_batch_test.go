package deploy

import "testing"

func TestNewPlanRejectsInvalidStrategyAndTarget(t *testing.T) {
	target := Target{ID:"t1", Name:"node", Serial:"serial-1"}
	if _, err := NewPlan("p1", "project", target, "artifact", "invalid"); err == nil { t.Fatal("unsupported strategy must fail") }
	if _, err := NewPlan("p1", "project", Target{}, "artifact", "rolling"); err == nil { t.Fatal("invalid target must fail") }
}
