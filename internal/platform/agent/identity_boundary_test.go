package agent

import "testing"

func TestAgentPlanLookupIsStable(t *testing.T) {
	if Plans["free"].ID != "free" || Plans["pro"].ID != "pro" {
		t.Fatal("expected built-in plans to remain registered")
	}
}
