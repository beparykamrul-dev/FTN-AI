package mesh

import "testing"

func TestSessionSubscribe(t *testing.T) {
 r:=NewSessionRegistry(); s:=r.Open("node"); if !r.Subscribe(s.ID,Heartbeat){t.Fatal("subscription failed")}
}
