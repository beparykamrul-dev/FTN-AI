package mesh

import "testing"

func TestSessionOpenTrimsNodeID(t *testing.T) {
 r:=NewSessionRegistry(); s:=r.Open("  node-1  "); if s.NodeID!="node-1" || s.ID=="" {t.Fatalf("unexpected session: %#v",s)}
}
