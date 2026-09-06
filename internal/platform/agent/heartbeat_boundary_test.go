package agent

import (
	"testing"
	"time"
)

func TestRegistryRejectsFutureHeartbeat(t *testing.T) {
	r := NewRegistry()
	state := r.Upsert(Heartbeat{AgentID: "a1", HostID: "h1", ObservedAt: time.Now().UTC().Add(10 * time.Minute)})
	if state.AgentID != "" {
		t.Fatalf("future heartbeat was accepted: %+v", state)
	}
}
