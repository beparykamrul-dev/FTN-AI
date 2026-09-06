package agent

import (
 "testing"
 "time"
)

func TestRegistryRejectsFarFutureHeartbeat(t *testing.T) {
 r:=NewRegistry(); s:=r.Upsert(Heartbeat{AgentID:"a",ObservedAt:time.Now().UTC().Add(6*time.Minute)})
 if s.AgentID!="" {t.Fatal("future heartbeat accepted")}
}
