package mesh

import (
 "testing"
 "time"
)

func TestHealthFreshBoundary(t *testing.T) {
 now:=time.Now().UTC(); r:=NewHealthRegistry(HealthPolicy{HeartbeatTimeout:time.Minute}); r.Observe("p",now,80)
 if !r.Fresh("p",now.Add(time.Minute)){t.Fatal("boundary freshness rejected")}
 if r.Fresh("p",now.Add(time.Minute+time.Nanosecond)){t.Fatal("stale peer reported fresh")}
}
