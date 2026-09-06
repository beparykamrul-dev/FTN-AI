package mesh

import "testing"

func TestSessionCloseRemovesSession(t *testing.T) {
 r:=NewSessionRegistry(); s:=r.Open("node"); r.Close(s.ID); if r.Heartbeat(s.ID){t.Fatal("closed session still active")}
}
