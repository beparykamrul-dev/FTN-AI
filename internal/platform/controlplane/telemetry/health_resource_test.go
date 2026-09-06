package telemetry

import (
 "testing"
 "time"
)

func TestEvaluateResourcePressure(t *testing.T) {
 now:=time.Now().UTC(); h:=Heartbeat{NodeID:"n",ObservedAt:now,CPUPercent:90,MemoryPercent:10,Healthy:true}
 got:=Evaluate(h,now,time.Minute)
 if got.Status!=HealthDegraded {t.Fatalf("got %s",got.Status)}
}
