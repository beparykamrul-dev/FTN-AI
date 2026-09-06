package mesh

import (
 "testing"
 "time"
)

func TestHealthObserveRejectsEmptyPeer(t *testing.T) {
 r:=NewHealthRegistry(DefaultHealthPolicy()); if got:=r.Observe("  ",time.Now(),50); got.PeerID!="" {t.Fatal("empty peer accepted")}
}
