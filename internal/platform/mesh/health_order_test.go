package mesh

import (
 "testing"
 "time"
)

func TestHealthEvaluateSortsPeerIDs(t *testing.T) {
 now:=time.Now().UTC(); r:=NewHealthRegistry(DefaultHealthPolicy()); r.Observe("z",now,1); r.Observe("a",now,1)
 got:=r.Evaluate(now)
 if len(got)!=2 || got[0].PeerID!="a" || got[1].PeerID!="z" {t.Fatalf("unexpected order: %#v",got)}
}
