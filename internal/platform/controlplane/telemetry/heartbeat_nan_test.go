package telemetry

import (
 "math"
 "testing"
 "time"
)

func TestHeartbeatRejectsNaNAndInf(t *testing.T) {
 h:=Heartbeat{NodeID:"n",ObservedAt:time.Now(),CPUPercent:math.NaN()}
 if h.Valid(){t.Fatal("NaN accepted")}
 h.CPUPercent=math.Inf(1)
 if h.Valid(){t.Fatal("Inf accepted")}
}
