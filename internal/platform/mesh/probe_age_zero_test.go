package mesh

import (
 "testing"
 "time"
)

func TestProbeAgeZeroObservedIsMaxDuration(t *testing.T) {
 if ProbeAge(time.Now(),time.Time{})!=time.Duration(1<<63-1){t.Fatal("zero observation not treated as oldest")}
}
