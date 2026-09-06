package telemetry

import (
 "testing"
 "time"
)

func TestEvaluateStaleHeartbeat(t *testing.T) {
 now:=time.Now().UTC(); h:=Heartbeat{NodeID:"n",ObservedAt:now.Add(-2*time.Minute),Healthy:true}
 got:=Evaluate(h,now,time.Minute)
 if got.Status!=HealthStale {t.Fatalf("got %s",got.Status)}
}
