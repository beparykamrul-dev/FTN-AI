package agent

import (
	"testing"
	"time"
)

func TestRegistryRejectsFarFutureHeartbeat(t *testing.T) {
	r := NewRegistry()
	if got := r.Upsert(Heartbeat{AgentID: "agent-1", ObservedAt: time.Now().Add(10 * time.Minute)}); got.AgentID != "" { t.Fatal("expected future heartbeat to be rejected") }
	if _, ok := r.Get("agent-1"); ok { t.Fatal("rejected heartbeat must not be stored") }
}

func TestRegistryListIsDeterministic(t *testing.T) {
	r := NewRegistry()
	r.Upsert(Heartbeat{AgentID: "b"})
	r.Upsert(Heartbeat{AgentID: "a"})
	items := r.List()
	if len(items) != 2 || items[0].AgentID != "a" || items[1].AgentID != "b" { t.Fatalf("unexpected order: %+v", items) }
}
