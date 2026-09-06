package mesh

import (
 "testing"
 "time"
)

func TestProbeAgeFutureObservationIsZero(t *testing.T) {
 now:=time.Now().UTC(); if ProbeAge(now,now.Add(time.Second))!=0 {t.Fatal("future probe produced negative age")}
}
