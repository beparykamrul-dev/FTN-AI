package agent

import (
 "testing"
 "time"
)

func TestRegistryUpsertNormalizesFields(t *testing.T) {
 r:=NewRegistry(); now:=time.Now(); s:=r.Upsert(Heartbeat{AgentID:" a ",HostID:" h ",Version:" v ",Status:" ONLINE ",ObservedAt:now})
 if s.AgentID!="a" || s.HostID!="h" || s.Version!="v" || s.Status!="online" {t.Fatalf("unexpected state: %#v",s)}
}
