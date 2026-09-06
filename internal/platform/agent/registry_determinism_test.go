package agent

import "testing"

func TestRegistryListIsDeterministic(t *testing.T) {
	r := NewRegistry()
	r.Upsert(Heartbeat{AgentID: "b", HostID: "h2"})
	r.Upsert(Heartbeat{AgentID: "a", HostID: "h1"})
	got := r.List()
	if len(got) != 2 || got[0].AgentID != "a" || got[1].AgentID != "b" {
		t.Fatalf("unexpected registry order: %+v", got)
	}
}
