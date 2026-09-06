package telemetry

import (
 "testing"
 "time"
)

func TestFreshAcceptsExactMaxAge(t *testing.T) {
 now:=time.Now().UTC(); h:=Heartbeat{NodeID:"n",ObservedAt:now.Add(-time.Minute)}
 if !Fresh(h,now,time.Minute){t.Fatal("exact max age rejected")}
}
