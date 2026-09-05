package observability

import (
	"testing"
	"time"
)

func TestRateHandlesCounterResetWithoutNegativeValues(t *testing.T) {
	got := Rate(TrafficSample{Interface:"eth0",Bytes:100,Packets:100}, TrafficSample{Interface:"eth0",Bytes:10,Packets:10}, time.Second)
	if got.BPS != 0 || got.PPS != 0 { t.Fatalf("counter reset must not produce negative rate: %#v", got) }
}

func TestRateRejectsNonPositiveElapsedAsZeroRate(t *testing.T) {
	got := Rate(TrafficSample{Interface:"eth0"}, TrafficSample{Interface:"eth0",Bytes:10}, 0)
	if got.BPS != 0 || got.PPS != 0 { t.Fatalf("non-positive elapsed must produce zero rate: %#v", got) }
}
