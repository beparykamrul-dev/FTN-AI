package agent

import (
	"context"
	"testing"
)

func TestUsageGateRejectsUnknownPlan(t *testing.T) {
	gate := NewUsageGate(PlansAsSlice())
	if err := gate.CheckAndConsume(context.Background(), "principal", "missing", 1); err == nil {
		t.Fatal("expected unknown plan error")
	}
}

func PlansAsSlice() []Plan {
	out := make([]Plan, 0, len(Plans))
	for _, p := range Plans {
		out = append(out, p)
	}
	return out
}
