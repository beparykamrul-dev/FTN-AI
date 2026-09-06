package agent

import "testing"

func TestRegistryListSortsAgentIDs(t *testing.T) {
 r:=NewRegistry(); r.Upsert(Heartbeat{AgentID:"z"}); r.Upsert(Heartbeat{AgentID:"a"}); got:=r.List(); if len(got)!=2 || got[0].AgentID!="a" || got[1].AgentID!="z" {t.Fatalf("unexpected order: %#v",got)}
}
