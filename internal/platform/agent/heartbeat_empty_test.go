package agent

import "testing"

func TestRegistryRejectsEmptyAgentID(t *testing.T) {
 r:=NewRegistry(); if s:=r.Upsert(Heartbeat{HostID:"h"}); s.AgentID!="" {t.Fatal("empty agent accepted")}
}
