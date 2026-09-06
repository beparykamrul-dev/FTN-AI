package agent

import "testing"

func TestAgentDefaultsExposeFreePlan(t *testing.T) {
	p, ok := Plans["free"]
	if !ok || p.ID != "free" {
		t.Fatal("free plan must remain available")
	}
}
