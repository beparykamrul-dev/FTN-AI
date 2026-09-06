package telemetry

import (
 "testing"
 "time"
)

func TestFreshRejectsFutureObservation(t *testing.T) {
 now:=time.Now().UTC()
 h:=Heartbeat{NodeID:"n",ObservedAt:now.Add(time.Second)}
 if Fresh(h,now,time.Minute){t.Fatal("future heartbeat accepted")}
}
